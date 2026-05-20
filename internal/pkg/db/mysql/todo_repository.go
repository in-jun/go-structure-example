package mysql

import (
	"context"
	"database/sql"

	"github.com/in-jun/go-structure-example/internal/app/todo"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) todo.Repository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Save(ctx context.Context, t *todo.Todo) error {
	query := `
		INSERT INTO todos (user_id, title, description, due_date)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, t.UserID, t.Title, t.Description, t.DueDate)
	if err != nil {
		return errors.Internal("Failed to create todo")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get todo ID")
	}

	t.ID = uint(id)
	return nil
}

func (r *todoRepository) Update(ctx context.Context, t *todo.Todo) error {
	query := `
		UPDATE todos
		SET title = ?, description = ?, due_date = ?, status = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		t.Title, t.Description, t.DueDate, string(t.Status), t.ID, t.UserID)
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
	result, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id = ?", id)
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

func (r *todoRepository) FindById(ctx context.Context, id uint) (*todo.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos WHERE id = ?
	`

	var t todo.Todo
	var status string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description,
		&status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get todo")
	}

	t.Status = todo.Status(status)
	return &t, nil
}

func (r *todoRepository) FindByUserId(ctx context.Context, userID uint, page, limit int) ([]*todo.Todo, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM todos WHERE user_id = ?", userID,
	).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	todos := make([]*todo.Todo, 0)
	for rows.Next() {
		var t todo.Todo
		var status string
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		t.Status = todo.Status(status)
		todos = append(todos, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating todo rows")
	}

	return todos, total, nil
}

func (r *todoRepository) UpdateStatus(ctx context.Context, id uint, status todo.Status) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE todos SET status = ? WHERE id = ?", string(status), id)
	if err != nil {
		return errors.Internal("Failed to update todo status")
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

func (r *todoRepository) FindByUserIdAndStatus(ctx context.Context, userID uint, status todo.Status, page, limit int) ([]*todo.Todo, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM todos WHERE user_id = ? AND status = ?", userID, string(status),
	).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos
		WHERE user_id = ? AND status = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, string(status), limit, (page-1)*limit)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	todos := make([]*todo.Todo, 0)
	for rows.Next() {
		var t todo.Todo
		var s string
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&s, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		t.Status = todo.Status(s)
		todos = append(todos, &t)
	}

	return todos, total, nil
}

func (r *todoRepository) FindUpcoming(ctx context.Context, userID uint, limit int) ([]*todo.Todo, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM todos
		WHERE user_id = ?
		  AND status = 'pending'
		  AND due_date > NOW()
		ORDER BY due_date ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, errors.Internal("Failed to get upcoming todos")
	}
	defer rows.Close()

	todos := make([]*todo.Todo, 0)
	for rows.Next() {
		var t todo.Todo
		var status string
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, errors.Internal("Failed to scan todo")
		}
		t.Status = todo.Status(status)
		todos = append(todos, &t)
	}

	return todos, nil
}
