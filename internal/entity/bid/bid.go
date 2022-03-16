package bid

import "errors"

const AuctionRef = "Auction"

type Bid struct {
	BidID      string `json:"bid_id"`
	ProductID  string `json:"product_id"`
	UserID     string `json:"user_id"`
	Amount     int64  `json:"amount"`
	PlacedTime int64  `json:"placed_time"`
	State      string `json:"string"`
}

type BidFirebaseRDB struct {
	ProductID    string `json:"product_id"`
	UserID       string `json:"user_id"`
	CurrentPrice int64  `json:"current_price"`
	BidderCount  int64  `json:"bidder_count"`
	Finished     bool   `json:"finished"`
}

const (
	StateBooked   = "book"
	StateLost     = "lost"
	StateReturned = "returned"
)

func (b *Bid) Validate() error {
	if b.UserID == "" {
		return errors.New("Invalid User")
	}
	if b.ProductID == "" {
		return errors.New("Invalid Product")
	}
	if b.Amount <= 0 {
		return errors.New("Invalid Amount")
	}
	return nil
}

func (b *Bid) ValidateMultiplierAmount() error {

	return nil
}

func (b *Bid) ValidateBidEligibility(balance, multiplier, highestBid, endTime int64) error {
	if b.Amount%multiplier != 0 {
		return errors.New("Invalid Multiplier Amount")
	}
	if b.Amount > balance {
		return errors.New("Invalid Remaining Wallet amount")
	}
	if b.Amount <= highestBid {
		return errors.New("Put Higher Amount, you poor guy")
	}
	if b.PlacedTime > endTime {
		return errors.New("Too late bro")
	}
	return nil
}

func GetListProductIDFromListBid(bids []Bid) map[string]bool {
	res := map[string]bool{}
	for _, bid := range bids {
		res[bid.ProductID] = true
	}
	return res
}
