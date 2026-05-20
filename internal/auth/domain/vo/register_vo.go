package vo

import (
	"errors"
	"net/mail"
)

var (
	errInvalidRegister         = errors.New("valid email, password (min 6 chars), and name are required")
	errInvalidRegisterEmail    = errors.New("invalid email format")
	errInvalidRegisterPassword = errors.New("password must be at least 6 characters")
)

type RegisterVO struct {
	Email    string
	Password string
	Name     string
}

func NewRegisterVO(email, password, name string) (*RegisterVO, error) {
	if email == "" || password == "" || name == "" {
		return nil, errInvalidRegister
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errInvalidRegisterEmail
	}
	if len(password) < 6 {
		return nil, errInvalidRegisterPassword
	}
	return &RegisterVO{Email: email, Password: password, Name: name}, nil
}
