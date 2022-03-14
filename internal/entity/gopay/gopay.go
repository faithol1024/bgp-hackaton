package gopay

import "errors"

type GopaySaldo struct {
	UserID       int64 `json:"user_gopay"`
	AmountIDR    int64 `json:"amount_idr"`
	AmountPoints int64 `json:"amount_points"`
}

func (g *GopaySaldo) Validate() error {
	if g.UserID == 0 {
		return errors.New("Invalid User")
	}
	return nil
}
