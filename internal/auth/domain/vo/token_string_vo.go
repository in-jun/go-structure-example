package vo

import "errors"

var errEmptyTokenString = errors.New("token string is required")

type TokenStringVO struct {
	Token string
}

func NewTokenStringVO(token string) (*TokenStringVO, error) {
	if token == "" {
		return nil, errEmptyTokenString
	}
	return &TokenStringVO{Token: token}, nil
}
