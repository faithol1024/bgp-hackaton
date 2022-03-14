package gopay

import (
	"context"
	"errors"

	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
)

type gopayResource interface {
	GetByUserID(ctx context.Context, userID int64) (gopay.GopaySaldo, error)
}

type UseCase struct {
	gopayRsc gopayResource
}

func New(gopayResource gopayResource) *UseCase {
	return &UseCase{
		gopayRsc: gopayResource,
	}
}

func (uc *UseCase) GetByUserID(ctx context.Context, userID int64) (gopay.GopaySaldo, error) {
	if userID == 0 {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(errors.New("invalid user_id"))
	}
	return uc.gopayRsc.GetByUserID(ctx, userID)
}