package vo

import "errors"

type LoginVO struct {
	Email    string
	Password string
}

func NewLoginVO(email, password string) (*LoginVO, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}
	return &LoginVO{Email: email, Password: password}, nil
}
