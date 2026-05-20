package vo

import (
	"errors"
	"time"
)

var (
	errUpdateTitleRequired = errors.New("title is required")
	errUpdateTitleTooLong  = errors.New("title must be 255 characters or less")
)

type UpdateTodoVO struct {
	Title       string
	Description string
	DueDate     time.Time
}

func NewUpdateTodoVO(title, description string, dueDate time.Time) (*UpdateTodoVO, error) {
	if title == "" {
		return nil, errUpdateTitleRequired
	}
	if len(title) > 255 {
		return nil, errUpdateTitleTooLong
	}
	return &UpdateTodoVO{Title: title, Description: description, DueDate: dueDate}, nil
}
