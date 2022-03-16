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
