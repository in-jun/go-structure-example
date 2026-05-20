package vo

import "errors"

var errInvalidUserID = errors.New("user ID must be greater than zero")

type UserIDVO struct {
	ID uint
}

func NewUserIDVO(id uint) (*UserIDVO, error) {
	if id == 0 {
		return nil, errInvalidUserID
	}
	return &UserIDVO{ID: id}, nil
}
