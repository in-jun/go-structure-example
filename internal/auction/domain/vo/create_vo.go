package vo

import (
	"errors"
	"time"
)

var errInvalidCreate = errors.New("title and positive start price are required")

type CreateVO struct {
	Title       string
	Description string
	StartPrice  int64
	EndTime     time.Time
}

func NewCreateVO(title, description string, startPrice int64, endTime time.Time) (*CreateVO, error) {
	if title == "" || startPrice <= 0 {
		return nil, errInvalidCreate
	}
	return &CreateVO{Title: title, Description: description, StartPrice: startPrice, EndTime: endTime}, nil
}
