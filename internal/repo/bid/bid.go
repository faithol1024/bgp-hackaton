package bid

import (
	"context"

	database "firebase.google.com/go/v4/db"
	"github.com/faithol1024/bgp-hackaton/internal/entity/bid"
	"github.com/tokopedia/tdk/go/log"
)

type Repo struct {
	frdb *database.Ref
}

func New(frdb *database.Ref) *Repo {
	return &Repo{
		frdb: frdb,
	}
}

func (r *Repo) PublishBidFRDB(ctx context.Context, bid bid.BidFirebaseRDB) error {
	err := r.frdb.Child(bid.ProductID).Set(ctx, bid)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *Repo) Bid(ctx context.Context, bidReq bid.Bid) (bid.Bid, error) {
	return bid.Bid{}, nil
}

func (r *Repo) AntiDoubleRequest(ctx context.Context, userID string) error {
	return nil
}
