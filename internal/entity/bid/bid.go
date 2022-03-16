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
	if b.Amount <= 0 {
		return errors.New("Invalid Amount")
	}
	return nil
}

func (b *Bid) ValidateMultiplierAmount() error {

	return nil
}

func (b *Bid) ValidateBidEligibility(balance, multiplier, highestBid int64) error {
	if b.Amount%multiplier != 0 {
		return errors.New("Invalid Multiplier Amount")
	}
	if b.Amount > balance {
		return errors.New("Invalid Remaining Wallet amount")
	}
	if b.Amount <= highestBid {
		return errors.New("Put Higher Amount, you poor guy")
	}
	return nil
}
