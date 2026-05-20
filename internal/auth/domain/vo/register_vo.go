package vo

import (
	"errors"
	"net/mail"
)

type RegisterVO struct {
	Email    string
	Password string
	Name     string
}

func NewRegisterVO(email, password, name string) (*RegisterVO, error) {
	if email == "" || password == "" || name == "" {
		return nil, errors.New("email, password, and name are required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	return &RegisterVO{Email: email, Password: password, Name: name}, nil
}
