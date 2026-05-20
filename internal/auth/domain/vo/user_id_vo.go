package vo

import (
	"errors"

	"github.com/google/uuid"
)

var errInvalidUserID = errors.New("user ID must be a valid UUID")

type UserIDVO struct {
	ID string
}

func NewUserIDVO(id string) (*UserIDVO, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errInvalidUserID
	}
	return &UserIDVO{ID: id}, nil
}
