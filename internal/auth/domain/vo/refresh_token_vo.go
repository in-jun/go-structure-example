package vo

import "errors"

var errEmptyRefreshToken = errors.New("refresh token is required")

type RefreshTokenVO struct {
	Token string
}

func NewRefreshTokenVO(token string) (*RefreshTokenVO, error) {
	if token == "" {
		return nil, errEmptyRefreshToken
	}
	return &RefreshTokenVO{Token: token}, nil
}
