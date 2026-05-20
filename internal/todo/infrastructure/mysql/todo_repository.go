package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

var _ domain.TodoRepository = (*todoRepository)(nil)

type todoRepository struct {
	db func(ctx context.Context) transaction.DBTX
}

func NewTodoRepository(db func(ctx context.Context) transaction.DBTX) domain.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Save(ctx context.Context, t *entity.Todo) error {
	query := "INSERT INTO todos (user_id, title, description, due_date) VALUES (?, ?, ?, ?)"
	result, err := r.db(ctx).ExecContext(ctx, query, t.UserID(), t.Title(), t.Description(), t.DueDate())
	if err != nil {
		return errors.Internal("Failed to create todo")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get todo ID")
	}
	t.SetID(uint(id))
	return nil
}

func (r *todoRepository) FindByID(ctx context.Context, id uint) (*entity.Todo, error) {
	query := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at FROM todos WHERE id = ?"
	var tid, userID uint
	var title, description, status string
	var dueDate, createdAt, updatedAt time.Time
	err := r.db(ctx).QueryRowContext(ctx, query, id).Scan(
		&tid, &userID, &title, &description, &status, &dueDate, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get todo")
	}
	return entity.ReconstructTodo(tid, userID, title, description, entity.Status(status), dueDate, createdAt, updatedAt), nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]*entity.Todo, int64, error) {
	var total int64
	if err := r.db(ctx).QueryRowContext(ctx, "SELECT COUNT(*) FROM todos WHERE user_id = ?", userID).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db(ctx).QueryContext(ctx, query, userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	todos := make([]*entity.Todo, 0)
	for rows.Next() {
		var tid, uid uint
		var title, description, status string
		var dueDate, createdAt, updatedAt time.Time
		if err := rows.Scan(&tid, &uid, &title, &description, &status, &dueDate, &createdAt, &updatedAt); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		todos = append(todos, entity.ReconstructTodo(tid, uid, title, description, entity.Status(status), dueDate, createdAt, updatedAt))
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating todo rows")
	}

	return todos, total, nil
}

func (r *todoRepository) Update(ctx context.Context, t *entity.Todo) error {
	query := "UPDATE todos SET title = ?, description = ?, due_date = ?, status = ? WHERE id = ? AND user_id = ?"
	result, err := r.db(ctx).ExecContext(ctx, query,
		t.Title(), t.Description(), t.DueDate(), string(t.Status()), t.ID(), t.UserID())
	if err != nil {
		return errors.Internal("Failed to update todo")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}
	if rows == 0 {
		return errors.NotFound("Todo not found")
	}
	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	result, err := r.db(ctx).ExecContext(ctx, "DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return errors.Internal("Failed to delete todo")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}
	if rows == 0 {
		return errors.NotFound("Todo not found")
	}
	return nil
}
