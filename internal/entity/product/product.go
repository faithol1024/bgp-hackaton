package product

import (
	"errors"

	ers "github.com/faithol1024/bgp-hackaton/lib/error"
)

type Product struct {
	ProductID    string `json:"product_id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	ImageURL     string `json:"image_url"`
	Description  string `json:"description"`
	StartBid     int64  `json:"start_bid"`
	MultipleBid  int64  `json:"multiple_bid"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	HighestBidID string `json:"highest_bid_id"`
	TotalBidder  int64  `json:"total_bidder"`
	Status       string `json:"status"`
}

const (
	StatusDone = "done"
	StatusNew  = "new"
)

func (p *Product) Validate() error {
	if p.UserID == "" {
		return ers.ErrorAddTrace(errors.New("Invalid User"))
	}
	if p.ProductID == "" {
		return ers.ErrorAddTrace(errors.New("Invalid Product"))
	}
	return nil
}
