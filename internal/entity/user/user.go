package user

import (
	"errors"
	"strings"
)

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

func (u *User) GetMaskedName() string {
	temp := strings.Fields(u.UserName)
	res := ""
	for _, t := range temp {
		word := ""
		for i := 0; i < len(t); i++ {
			if i == 0 {
				word += string(t[i])
			}
			word += "*"
		}
		res += word + " "
	}
	res = strings.TrimSpace(res)
	return res
}
