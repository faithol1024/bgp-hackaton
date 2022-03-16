package user

import (
	"context"
	"errors"
	"github.com/faithol1024/bgp-hackaton/internal/usecase/gopay"
	"github.com/faithol1024/bgp-hackaton/lib/util"

	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
)

type userResource interface {
	GetByID(ctx context.Context, userID string) (user.User, error)
	Create(ctx context.Context, user user.User) (user.User, error)
}

type UseCase struct {
	userRsc userResource
	gopayUC *gopay.UseCase
}

func New(userResource userResource, gopayUC *gopay.UseCase) *UseCase {
	return &UseCase{
		userRsc: userResource,
		gopayUC: gopayUC,
	}
}

func (uc *UseCase) GetByID(ctx context.Context, userID string) (user.User, error) {
	if userID == "" {
		return user.User{}, ers.ErrorAddTrace(errors.New("invalid user_id"))
	}
	return uc.userRsc.GetByID(ctx, userID)
}

func (uc *UseCase) Create(ctx context.Context, req user.User) (user.User, error) {
	req.UserID = util.GetStringUUID()
	res, err := uc.userRsc.Create(ctx, req)
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}
	_, err = uc.gopayUC.Create(ctx, res.UserID)
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}
	return res, err
}
