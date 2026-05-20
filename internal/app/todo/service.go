package todo

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
)

type Service interface {
	Create(ctx context.Context, userID uint, req CreateTodoRequest) (*TodoResponse, error)
	GetList(ctx context.Context, userID uint, page, limit int) (*TodoListResponse, error)
	GetByID(ctx context.Context, userID, todoID uint) (*TodoResponse, error)
	Update(ctx context.Context, userID, todoID uint, req UpdateTodoRequest) error
	UpdateStatus(ctx context.Context, userID, todoID uint, req UpdateTodoStatusRequest) error
	Delete(ctx context.Context, userID, todoID uint) error
}

var _ Service = (*service)(nil)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uint, req CreateTodoRequest) (*TodoResponse, error) {
	t := &Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      StatusPending,
	}

	if err := s.repo.Save(ctx, t); err != nil {
		return nil, err
	}

	return toResponse(t), nil
}

func (s *service) GetList(ctx context.Context, userID uint, page, limit int) (*TodoListResponse, error) {
	todos, total, err := s.repo.FindByUserId(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	response := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		response = append(response, *toResponse(t))
	}

	return &TodoListResponse{Todos: response, Total: total}, nil
}

func (s *service) GetByID(ctx context.Context, userID, todoID uint) (*TodoResponse, error) {
	t, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("Todo not found")
	}
	if t.UserID != userID {
		return nil, errors.Forbidden("Not authorized to access this todo")
	}

	return toResponse(t), nil
}

func (s *service) Update(ctx context.Context, userID, todoID uint, req UpdateTodoRequest) error {
	t, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.NotFound("Todo not found")
	}
	if t.UserID != userID {
		return errors.Forbidden("Not authorized to update this todo")
	}

	t.Title = req.Title
	t.Description = req.Description
	t.DueDate = req.DueDate

	return s.repo.Update(ctx, t)
}

func (s *service) UpdateStatus(ctx context.Context, userID, todoID uint, req UpdateTodoStatusRequest) error {
	t, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.NotFound("Todo not found")
	}
	if t.UserID != userID {
		return errors.Forbidden("Not authorized to update this todo")
	}

	return s.repo.UpdateStatus(ctx, todoID, req.Status)
}

func (s *service) Delete(ctx context.Context, userID, todoID uint) error {
	t, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.NotFound("Todo not found")
	}
	if t.UserID != userID {
		return errors.Forbidden("Not authorized to delete this todo")
	}

	return s.repo.Delete(ctx, todoID)
}

func toResponse(t *Todo) *TodoResponse {
	return &TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
