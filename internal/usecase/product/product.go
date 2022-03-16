package product

import (
	"context"
	"errors"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	"github.com/faithol1024/bgp-hackaton/internal/entity/product"
	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/faithol1024/bgp-hackaton/lib/util"
)

type productResource interface {
	Create(ctx context.Context, product product.Product) error
	GetByID(ctx context.Context, ID string) (product.Product, error)
	GetAll(ctx context.Context) ([]product.Product, error)
	GetAllBySeller(ctx context.Context, userID string) ([]product.Product, error)
	GetAllByBuyer(ctx context.Context, userID string) ([]product.Product, error)
}

type bidResource interface {
	Bid(ctx context.Context, bid bid.Bid) (bid.Bid, error)
	// GetBidByProduct(ctx context.Context, productID string) (bid.Bid, error)
	AntiDoubleRequest(ctx context.Context, userID string) error
	// ReleaseAntiDoubleRequest(ctx context.Context, userID string) error
}

type gopayResource interface {
	GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error)
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

func (uc *UseCase) Create(ctx context.Context, product product.Product) error {
	product.ProductID = util.GetStringUUID()

	err := product.Validate()
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	err = uc.productRsc.Create(ctx, product)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (uc *UseCase) GetByID(ctx context.Context, ID string) (product.Product, error) {
	productRes, err := uc.productRsc.GetByID(ctx, ID)
	if err != nil {
		return product.Product{}, ers.ErrorAddTrace(err)
	}

	err = productRes.Validate()
	if err != nil {
		return product.Product{}, ers.ErrorAddTrace(err)
	}

	return productRes, nil
}

func (uc *UseCase) GetAll(ctx context.Context, userID string, role string) (products []product.Product, err error) {
	switch role {
	case user.RoleBuyer:
		products, err = uc.productRsc.GetAllByBuyer(ctx, userID)
	case user.RoleSeller:
		products, err = uc.productRsc.GetAllBySeller(ctx, userID)
	default:
		products, err = uc.productRsc.GetAll(ctx)
	}
	if err != nil {
		return []product.Product{}, ers.ErrorAddTrace(err)
	}

	if len(products) == 0 {
		return []product.Product{}, ers.ErrorAddTrace(errors.New("No products available"))
	}

	return products, nil
}

func (uc *UseCase) Bid(ctx context.Context, bidReq bid.Bid) (bid.Bid, error) {
	err := uc.bidRsc.AntiDoubleRequest(ctx, bidReq.UserID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}
	// defer uc.bidRsc.ReleaseAntiDoubleRequest(ctx, bidReq.UserID)

	gopay, err := uc.gopayRsc.GetByUserID(ctx, bidReq.UserID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	product, err := uc.productRsc.GetByID(ctx, bidReq.ProductID)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	// highestBid, err :=

	err = bidReq.ValidateBidEligibility(gopay.AmountIDR, product.MultipleBid, 0)
	if err != nil {
		return bid.Bid{}, ers.ErrorAddTrace(err)
	}

	return bid.Bid{}, nil
}
