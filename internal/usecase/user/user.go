package user

import (
	"context"
	"errors"

	"github.com/faithol1024/bgp-hackaton/internal/entity/user"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
)

type userResource interface {
	GetByUserID(ctx context.Context, userID string) (user.User, error)
	Create(ctx context.Context, user user.User) (user.User, error)
}

type UseCase struct {
	userRsc userResource
}

func New(userResource userResource) *UseCase {
	return &UseCase{
		userRsc: userResource,
	}
}

func (uc *UseCase) GetByUserID(ctx context.Context, userID string) (user.User, error) {
	if userID == "" {
		return user.User{}, ers.ErrorAddTrace(errors.New("invalid user_id"))
	}
	return uc.userRsc.GetByUserID(ctx, userID)
}

func (uc *UseCase) Create(ctx context.Context, req user.User) (user.User, error) {
	res, err := uc.userRsc.Create(ctx, req)
	if err != nil {
		return user.User{}, ers.ErrorAddTrace(err)
	}
	return res, err
}
