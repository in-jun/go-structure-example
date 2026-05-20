package vo

import (
	"errors"
	"time"
)

type CreateTodoVO struct {
	Title       string
	Description string
	DueDate     time.Time
}

func NewCreateTodoVO(title, description string, dueDate time.Time) (*CreateTodoVO, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if len(title) > 255 {
		return nil, errors.New("title must be 255 characters or less")
	}
	return &CreateTodoVO{Title: title, Description: description, DueDate: dueDate}, nil
}
