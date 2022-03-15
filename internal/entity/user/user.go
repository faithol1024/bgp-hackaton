package user

import "errors"

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("Invalid Name")
	}
	if u.Email == "" {
		return errors.New("Invalid Email")
	}
	return nil
}
