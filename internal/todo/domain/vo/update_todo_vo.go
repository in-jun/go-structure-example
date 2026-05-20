package vo

import (
	"errors"
	"time"
)

type UpdateTodoVO struct {
	Title       string
	Description string
	DueDate     time.Time
}

func NewUpdateTodoVO(title, description string, dueDate time.Time) (*UpdateTodoVO, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if len(title) > 255 {
		return nil, errors.New("title must be 255 characters or less")
	}
	return &UpdateTodoVO{Title: title, Description: description, DueDate: dueDate}, nil
}
