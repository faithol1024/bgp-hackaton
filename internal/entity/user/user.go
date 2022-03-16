package user

import "errors"

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

const (
	RoleSeller = "seller"
	RoleBuyer  = "buyer"
)

func (u *User) Validate() error {
	if u.UserName == "" {
		return errors.New("Invalid UserName")
	}
	if u.Email == "" {
		return errors.New("Invalid Email")
	}
	return nil
}
