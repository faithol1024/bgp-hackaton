package product

import (
	"context"

	"github.com/faithol1024/bgp-hackaton/internal/entity/product"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
	"github.com/faithol1024/bgp-hackaton/lib/util"
)

type productResource interface {
	Create(ctx context.Context, product product.Product) error
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
