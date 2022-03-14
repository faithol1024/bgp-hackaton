package bid

import (
	"context"

	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
)

type bidResource interface {
	PublishBidFRDB(ctx context.Context, bid bid.Bid) error
}

type UseCase struct {
	bidRsc bidResource
}

func New(bidResource bidResource) *UseCase {
	return &UseCase{
		bidRsc: bidResource,
	}
}

func (uc *UseCase) UpdateBidFRDB(ctx context.Context, bid bid.Bid) error {
	return uc.bidRsc.PublishBidFRDB(ctx, bid)
}
