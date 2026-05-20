package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

var (
	errInvalidTodo            = errors.New("user ID and title are required")
	errInvalidReconstructTodo = errors.New("id, user ID, and title are required")
)

type Todo struct {
	id          string
	userID      string
	title       string
	description string
	status      Status
	dueDate     time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTodo(userID, title, description string, dueDate time.Time) (*Todo, error) {
	if userID == "" || title == "" {
		return nil, errInvalidTodo
	}
	now := time.Now()
	return &Todo{
		id:          uuid.New().String(),
		userID:      userID,
		title:       title,
		description: description,
		status:      StatusPending,
		dueDate:     dueDate,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructTodo(id, userID, title, description string, status Status, dueDate, createdAt, updatedAt time.Time) (*Todo, error) {
	if id == "" || userID == "" || title == "" {
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

func (t *Todo) ID() string           { return t.id }
func (t *Todo) UserID() string       { return t.userID }
func (t *Todo) Title() string        { return t.title }
func (t *Todo) Description() string  { return t.description }
func (t *Todo) Status() Status       { return t.status }
func (t *Todo) DueDate() time.Time   { return t.dueDate }
func (t *Todo) CreatedAt() time.Time { return t.createdAt }
func (t *Todo) UpdatedAt() time.Time { return t.updatedAt }

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
