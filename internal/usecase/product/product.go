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
	Update(ctx context.Context, req productEntity.Product) error
	GetByID(ctx context.Context, ID string) (productEntity.Product, error)
	GetAll(ctx context.Context) ([]productEntity.Product, error)
	GetAllBySeller(ctx context.Context, userID string) ([]productEntity.Product, error)
	GetAllByBuyer(ctx context.Context, userProductIDs map[string]bool) ([]productEntity.Product, error)
}

type bidResource interface {
	Bid(ctx context.Context, bid bid.Bid, product productEntity.Product) (int64, error)
	GetHighestBidAmountByProduct(ctx context.Context, product productEntity.Product) (int64, error)
	GetHighestBidAmountByProductDB(ctx context.Context, productID string) (bid.Bid, error)
	AntiDoubleRequest(ctx context.Context, userID string) error
	ReleaseAntiDoubleRequest(ctx context.Context, userID string) error
	GetAllBidByUserID(ctx context.Context, userID string) ([]bid.Bid, error)
	SetHighestBidAmountByProductDB(ctx context.Context, bid bid.Bid) error
	PublishBidFRDB(ctx context.Context, bid bid.BidFirebaseRDB) error
	ReleaseBookedSaldo(ctx context.Context, productID string) error
}

type gopayResource interface {
	GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error)
	BookSaldo(ctx context.Context, userID, bidID string, amount int64) error
}

type userResource interface {
	GetByID(ctx context.Context, userID string) (user.User, error)
}

type UseCase struct {
	productRsc productResource
	bidRsc     bidResource
	gopayRsc   gopayResource
	userRsc    userResource
}

func New(productRsc productResource, bidRsc bidResource, gopayRsc gopayResource, userRsc userResource) *UseCase {
	return &UseCase{
		productRsc: productRsc,
		bidRsc:     bidRsc,
		gopayRsc:   gopayRsc,
		userRsc:    userRsc,
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
	//bidRes, err := uc.bidRsc.GetHighestBidAmountByProductDB(ctx, productRes.ProductID)
	//if err != nil {
	//	return productEntity.Product{}, ers.ErrorAddTrace(err)
	//}
	userRes, err := uc.userRsc.GetByID(ctx, productRes.UserID)
	if err != nil {
		return productEntity.Product{}, ers.ErrorAddTrace(err)
	}
	if time.Now().Unix() >= productRes.EndTime {
		productRes.UserName = userRes.GetMaskedName()
	} else {
		productRes.UserName = userRes.UserName
	}

	err = uc.FinishBid(ctx, productRes, userRes.UserID)
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

func (uc *UseCase) FinishBid(ctx context.Context, product productEntity.Product, userID string) error {
	if time.Now().Unix() >= product.EndTime {
		//err := uc.productRsc.Update(ctx, product)
		highestBid, err := uc.bidRsc.GetHighestBidAmountByProduct(ctx, product)
		if err != nil {
			return ers.ErrorAddTrace(err)
		}
		//if err != nil {
		//	return ers.ErrorAddTrace(err)
		//}
		go uc.bidRsc.PublishBidFRDB(ctx, bid.BidFirebaseRDB{
			ProductID:    product.ProductID,
			UserID:       userID,
			CurrentPrice: highestBid,
			BidderCount:  product.TotalBidder,
		})
	}
	return nil
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
	bidReq.State = bid.StateBooked
	err = bidReq.ValidateBidEligibility(gopay.AmountIDR, product.MultipleBid, highestBid, product.EndTime)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	go uc.gopayRsc.BookSaldo(ctx, bidReq.UserID, bidReq.BidID, bidReq.Amount)

	count, err := uc.bidRsc.Bid(ctx, bidReq, product)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	go uc.bidRsc.PublishBidFRDB(ctx, bid.BidFirebaseRDB{
		ProductID:    bidReq.ProductID,
		UserID:       bidReq.UserID,
		CurrentPrice: bidReq.Amount,
		BidderCount:  count,
	})

	go uc.bidRsc.ReleaseBookedSaldo(ctx, bidReq.ProductID)

	go uc.bidRsc.SetHighestBidAmountByProductDB(ctx, bidReq)

	return bidReq, nil
}
