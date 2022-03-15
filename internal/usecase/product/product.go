package product

import (
	"context"
	"errors"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
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
	Bid(ctx context.Context, bid bid.Bid) (bid.Bid, error)
}

type UseCase struct {
	productRsc productResource
}

func New(productRsc productResource) *UseCase {
	return &UseCase{
		productRsc: productRsc,
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
	return bid.Bid{}, nil
}
