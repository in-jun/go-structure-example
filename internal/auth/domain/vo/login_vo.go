package vo

import "errors"

var errInvalidLogin = errors.New("email and password are required")

type LoginVO struct {
	Email    string
	Password string
}

func NewLoginVO(email, password string) (*LoginVO, error) {
	if email == "" || password == "" {
		return nil, errInvalidLogin
	}
	return &LoginVO{Email: email, Password: password}, nil
}
