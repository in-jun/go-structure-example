package pg

import (
	"context"
	"database/sql"
	stderrors "errors"
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
	query := `INSERT INTO todos (id, user_id, title, description, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db(ctx).ExecContext(ctx, query,
		t.ID(), t.UserID(), t.Title(), t.Description(), string(t.Status()), t.DueDate(), t.CreatedAt(), t.UpdatedAt())
	if err != nil {
		return errors.Internal("Failed to create todo")
	}
	return nil
}

func (r *todoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
	query := "SELECT id, user_id, title, description, status, due_date, created_at, updated_at FROM todos WHERE id = $1"
	var tid, userID, title, description, status string
	var dueDate, createdAt, updatedAt time.Time
	err := r.db(ctx).QueryRowContext(ctx, query, id).Scan(
		&tid, &userID, &title, &description, &status, &dueDate, &createdAt, &updatedAt,
	)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get todo")
	}
	t, err := entity.ReconstructTodo(tid, userID, title, description, entity.Status(status), dueDate, createdAt, updatedAt)
	if err != nil {
		return nil, errors.Internal("Failed to reconstruct todo")
	}
	return t, nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]*entity.Todo, int64, error) {
	var total int64
	if err := r.db(ctx).QueryRowContext(ctx, "SELECT COUNT(*) FROM todos WHERE user_id = $1", userID).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db(ctx).QueryContext(ctx, query, userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	todos := make([]*entity.Todo, 0)
	for rows.Next() {
		var tid, uid, title, description, status string
		var dueDate, createdAt, updatedAt time.Time
		if err := rows.Scan(&tid, &uid, &title, &description, &status, &dueDate, &createdAt, &updatedAt); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		t, err := entity.ReconstructTodo(tid, uid, title, description, entity.Status(status), dueDate, createdAt, updatedAt)
		if err != nil {
			return nil, 0, errors.Internal("Failed to reconstruct todo")
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating todo rows")
	}

	return todos, total, nil
}

func (r *todoRepository) Update(ctx context.Context, t *entity.Todo) error {
	query := `UPDATE todos SET title = $1, description = $2, due_date = $3, status = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7`
	result, err := r.db(ctx).ExecContext(ctx, query,
		t.Title(), t.Description(), t.DueDate(), string(t.Status()), t.UpdatedAt(), t.ID(), t.UserID())
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

func (r *todoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db(ctx).ExecContext(ctx, "DELETE FROM todos WHERE id = $1", id)
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
