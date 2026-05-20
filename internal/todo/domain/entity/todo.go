package entity

import (
	"errors"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

var (
	errInvalidTodo          = errors.New("user ID and title are required")
	errInvalidReconstructTodo = errors.New("id, user ID, and title are required")
)

type Todo struct {
	id          uint
	userID      uint
	title       string
	description string
	status      Status
	dueDate     time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTodo(userID uint, title, description string, dueDate time.Time) (*Todo, error) {
	if userID == 0 || title == "" {
		return nil, errInvalidTodo
	}
	now := time.Now()
	return &Todo{
		userID:      userID,
		title:       title,
		description: description,
		status:      StatusPending,
		dueDate:     dueDate,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructTodo(id, userID uint, title, description string, status Status, dueDate, createdAt, updatedAt time.Time) (*Todo, error) {
	if id == 0 || userID == 0 || title == "" {
		return nil, errInvalidReconstructTodo
	}
	return &Todo{
		id:          id,
		userID:      userID,
		title:       title,
		description: description,
		status:      status,
		dueDate:     dueDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

func (t *Todo) ID() uint             { return t.id }
func (t *Todo) UserID() uint         { return t.userID }
func (t *Todo) Title() string        { return t.title }
func (t *Todo) Description() string  { return t.description }
func (t *Todo) Status() Status       { return t.status }
func (t *Todo) DueDate() time.Time   { return t.dueDate }
func (t *Todo) CreatedAt() time.Time { return t.createdAt }
func (t *Todo) UpdatedAt() time.Time { return t.updatedAt }

func (t *Todo) SetID(id uint) { t.id = id }

func (t *Todo) Update(title, description string, dueDate time.Time) {
	t.title = title
	t.description = description
	t.dueDate = dueDate
	t.updatedAt = time.Now()
}

func (t *Todo) SetStatus(status Status) {
	t.status = status
	t.updatedAt = time.Now()
}
