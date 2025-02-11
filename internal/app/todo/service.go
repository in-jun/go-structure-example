package todo

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID uint, req CreateTodoRequest) (*TodoResponse, error) {
	todo := &Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      "pending",
	}

	if err := s.repo.Save(ctx, todo); err != nil {
		return nil, err
	}

	return &TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}

func (s *Service) GetList(ctx context.Context, userID uint, page, limit int) (*TodoListResponse, error) {
	todos, total, err := s.repo.FindByUserId(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	var response []TodoResponse
	for _, todo := range todos {
		response = append(response, TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			DueDate:     todo.DueDate,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	}

	return &TodoListResponse{
		Todos: response,
		Total: total,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, userID, todoID uint) (*TodoResponse, error) {
	todo, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, errors.NotFound("Todo not found")
	}

	if todo.UserID != userID {
		return nil, errors.Forbidden("Not authorized to access this todo")
	}

	return &TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		DueDate:     todo.DueDate,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, userID, todoID uint, req UpdateTodoRequest) error {
	todo, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}

	if todo == nil {
		return errors.NotFound("Todo not found")
	}

	if todo.UserID != userID {
		return errors.Forbidden("Not authorized to update this todo")
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.DueDate = req.DueDate

	return s.repo.Update(ctx, todo)
}

func (s *Service) UpdateStatus(ctx context.Context, userID, todoID uint, req UpdateTodoStatusRequest) error {
	todo, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}

	if todo == nil {
		return errors.NotFound("Todo not found")
	}

	if todo.UserID != userID {
		return errors.Forbidden("Not authorized to update this todo")
	}

	return s.repo.UpdateStatus(ctx, todoID, req.Status)
}

func (s *Service) Delete(ctx context.Context, userID, todoID uint) error {
	todo, err := s.repo.FindById(ctx, todoID)
	if err != nil {
		return err
	}

	if todo == nil {
		return errors.NotFound("Todo not found")
	}

	if todo.UserID != userID {
		return errors.Forbidden("Not authorized to delete this todo")
	}

	return s.repo.Delete(ctx, todoID)
}
