package product

import (
	"context"
	"errors"
	"time"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	productEntity "github.com/faithol1024/bgp-hackaton/internal/entity/product"
	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/faithol1024/bgp-hackaton/lib/util"
)

type productResource interface {
	Create(ctx context.Context, product productEntity.Product) error
	GetByID(ctx context.Context, ID string) (productEntity.Product, error)
	GetAll(ctx context.Context) ([]productEntity.Product, error)
	GetAllBySeller(ctx context.Context, userID string) ([]productEntity.Product, error)
	GetAllByBuyer(ctx context.Context, userProductIDs map[string]bool) ([]productEntity.Product, error)
}

type bidResource interface {
	Bid(ctx context.Context, bid bid.Bid, product productEntity.Product) (int64, error)
	GetHighestBidAmountByProduct(ctx context.Context, product productEntity.Product) (int64, error)
	AntiDoubleRequest(ctx context.Context, userID string) error
	ReleaseAntiDoubleRequest(ctx context.Context, userID string) error
	GetAllBidByUserID(ctx context.Context, userID string) ([]bid.Bid, error)
	SetHighestBidAmountByProductDB(ctx context.Context, bid bid.Bid) error
	PublishBidFRDB(ctx context.Context, bid bid.BidFirebaseRDB) error
}

type gopayResource interface {
	GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error)
	BookSaldo(ctx context.Context, userID, bidID string, amount int64) error
}

type UseCase struct {
	productRsc productResource
	bidRsc     bidResource
	gopayRsc   gopayResource
}

func New(productRsc productResource, bidRsc bidResource, gopayRsc gopayResource) *UseCase {
	return &UseCase{
		productRsc: productRsc,
		bidRsc:     bidRsc,
		gopayRsc:   gopayRsc,
	}
}

func (uc *UseCase) Create(ctx context.Context, product productEntity.Product) (productEntity.Product, error) {
	product.ProductID = util.GetStringUUID()

	product.StartTime = time.Now().Unix()
	product.Status = productEntity.StatusNew
	product.ProductID = util.GetStringUUID()
	err := uc.productRsc.Create(ctx, product)
	if err != nil {
		return product, ers.ErrorAddTrace(err)
	}

	return product, nil
}

func (uc *UseCase) GetByID(ctx context.Context, ID string) (productEntity.Product, error) {
	productRes, err := uc.productRsc.GetByID(ctx, ID)
	if err != nil {
		return productEntity.Product{}, ers.ErrorAddTrace(err)
	}

	err = productRes.Validate()
	if err != nil {
		return productEntity.Product{}, ers.ErrorAddTrace(err)
	}

	return productRes, nil
}

func (uc *UseCase) GetAll(ctx context.Context, userID string, role string) (products []productEntity.Product, err error) {
	switch role {
	case user.RoleBuyer:
		allBidByUser, err := uc.bidRsc.GetAllBidByUserID(ctx, userID)
		if err != nil {
			return []productEntity.Product{}, ers.ErrorAddTrace(err)
		}
		products, err = uc.productRsc.GetAllByBuyer(ctx, bid.GetListProductIDFromListBid(allBidByUser))
	case user.RoleSeller:
		products, err = uc.productRsc.GetAllBySeller(ctx, userID)
	default:
		products, err = uc.productRsc.GetAll(ctx)
	}
	if err != nil {
		return []productEntity.Product{}, ers.ErrorAddTrace(err)
	}

	if len(products) == 0 {
		return []productEntity.Product{}, ers.ErrorAddTrace(errors.New("No products available"))
	}

	return products, nil
}

func (uc *UseCase) Bid(ctx context.Context, bidReq bid.Bid) (bid.Bid, error) {
	err := uc.bidRsc.AntiDoubleRequest(ctx, bidReq.UserID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}
	defer uc.bidRsc.ReleaseAntiDoubleRequest(ctx, bidReq.UserID)

	gopay, err := uc.gopayRsc.GetByUserID(ctx, bidReq.UserID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	product, err := uc.productRsc.GetByID(ctx, bidReq.ProductID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	highestBid, err := uc.bidRsc.GetHighestBidAmountByProduct(ctx, product)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	bidReq.BidID = util.GetStringUUID()
	bidReq.PlacedTime = time.Now().Unix()
	err = bidReq.ValidateBidEligibility(gopay.AmountIDR, product.MultipleBid, highestBid)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	go uc.gopayRsc.BookSaldo(ctx, bidReq.UserID, bidReq.BidID, bidReq.Amount)

	count, err := uc.bidRsc.Bid(ctx, bidReq, product)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	go uc.bidRsc.PublishBidFRDB(ctx, bid.BidFirebaseRDB{
		ProductID:    bidReq.BidID,
		UserID:       bidReq.UserID,
		CurrentPrice: bidReq.Amount,
		BidderCount:  count,
	})

	go uc.bidRsc.SetHighestBidAmountByProductDB(ctx, bidReq)

	return bidReq, nil
}
