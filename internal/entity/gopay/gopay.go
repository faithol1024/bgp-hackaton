package gopay

import "errors"

type GopaySaldo struct {
	UserID    string `json:"user_id"`
	AmountIDR int64  `json:"amount_idr"`
}

func (g *GopaySaldo) Validate() error {
	if g.UserID == "" {
		return errors.New("Invalid User")
	}
	return nil
}
