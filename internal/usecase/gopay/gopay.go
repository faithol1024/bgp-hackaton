package gopay

import (
	"context"
	"errors"
	"github.com/faithol1024/bgp-hackaton/internal/entity/gopay"
	ers "github.com/faithol1024/bgp-hackaton/lib/error"
)

type gopayResource interface {
	GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error)
	GetAllHistoryByUserID(ctx context.Context, userID string) ([]gopay.GopayHistory, error)
	Create(ctx context.Context, req gopay.GopaySaldo) (gopay.GopaySaldo, error)
	CreateHistory(ctx context.Context, req gopay.GopayHistory) (gopay.GopayHistory, error)
}

type UseCase struct {
	gopayRsc gopayResource
}

func New(gopayResource gopayResource) *UseCase {
	return &UseCase{
		gopayRsc: gopayResource,
	}
}

func (uc *UseCase) GetByUserID(ctx context.Context, userID string) (gopay.GopaySaldo, error) {
	if userID == "" {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(errors.New("invalid user_id"))
	}
	return uc.gopayRsc.GetByUserID(ctx, userID)
}

func (uc *UseCase) GetAllHistoryByUserID(ctx context.Context, userID string) ([]gopay.GopayHistory, error) {
	if userID == "" {
		return nil, ers.ErrorAddTrace(errors.New("invalid user_id"))
	}
	return uc.gopayRsc.GetAllHistoryByUserID(ctx, userID)
}

func (uc *UseCase) Create(ctx context.Context, userID string) (gopay.GopaySaldo, error) {
	//1jt by default
	req := gopay.GopaySaldo{
		UserID:    userID,
		AmountIDR: 1000000,
	}

	err := req.Validate()
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}

	res, err := uc.gopayRsc.Create(ctx, req)
	if err != nil {
		return gopay.GopaySaldo{}, ers.ErrorAddTrace(err)
	}

	return res, nil
}
