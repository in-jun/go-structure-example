package vo

import (
	"errors"
	"time"
)

var (
	errTitleRequired = errors.New("title is required")
	errTitleTooLong  = errors.New("title must be 255 characters or less")
)

type CreateTodoVO struct {
	Title       string
	Description string
	DueDate     time.Time
}

func NewCreateTodoVO(title, description string, dueDate time.Time) (*CreateTodoVO, error) {
	if title == "" {
		return nil, errTitleRequired
	}
	if len(title) > 255 {
		return nil, errTitleTooLong
	}
	return &CreateTodoVO{Title: title, Description: description, DueDate: dueDate}, nil
}
