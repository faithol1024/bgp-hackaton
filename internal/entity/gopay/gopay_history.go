package gopay

import "errors"

type GopayHistory struct {
	GopayHistoryID string `json:"gopay_history_id"`
	UserID         string `json:"user_id"`
	GopayID        string `json:"gopay_id"`
	AmountIDR      int64  `json:"amount_idr"`
	BidID          string `json:"bid_id"`
}

func (g *GopayHistory) Validate() error {
	if g.UserID == "" {
		return errors.New("Invalid User")
	}
	return nil
}
