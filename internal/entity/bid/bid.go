package bid

import "errors"

const AuctionRef = "Auction"

type Bid struct {
	BidID      string `json:"bid_id"`
	ProductID  string `json:"product_id"`
	UserID     string `json:"user_id"`
	Amount     int64  `json:"amount"`
	PlacedTime int64  `json:"placed_time"`
}

func (b *Bid) Validate() error {
	if b.UserID == "" {
		return errors.New("Invalid User")
	}
	if b.ProductID == "" {
		return errors.New("Invalid Product")
	}
	return nil
}
