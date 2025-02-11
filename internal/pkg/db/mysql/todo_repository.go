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

func (r *todoRepository) Save(ctx context.Context, todo *todo.Todo) error {
	query := `
       INSERT INTO todos (user_id, title, description, due_date)
       VALUES (?, ?, ?, ?)
   `

	result, err := r.db.ExecContext(ctx, query,
		todo.UserID, todo.Title, todo.Description, todo.DueDate)
	if err != nil {
		return errors.Internal("Failed to create todo")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get todo ID")
	}

	todo.ID = uint(id)
	return nil
}

func (r *todoRepository) Update(ctx context.Context, todo *todo.Todo) error {
	query := `
       UPDATE todos 
       SET title = ?, description = ?, due_date = ?, status = ?
       WHERE id = ? AND user_id = ?
   `

	result, err := r.db.ExecContext(ctx, query,
		todo.Title, todo.Description, todo.DueDate, todo.Status, todo.ID, todo.UserID)
	if err != nil {
		return errors.Internal("Failed to update todo")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("Todo not found or not owned by user")
	}

	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM todos WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description,
		&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get todo")
	}

	return &t, nil
}

func (r *todoRepository) FindByUserId(ctx context.Context, userId uint, page, limit int) ([]todo.Todo, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM todos WHERE user_id = ?"
	if err := r.db.QueryRowContext(ctx, countQuery, userId).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	// Get todos with pagination
	query := `
       SELECT id, user_id, title, description, status, due_date, created_at, updated_at
       FROM todos 
       WHERE user_id = ?
       ORDER BY created_at DESC
       LIMIT ? OFFSET ?
   `

	rows, err := r.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	var todos []todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		todos = append(todos, t)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating todo rows")
	}

	return todos, total, nil
}

func (r *todoRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	query := "UPDATE todos SET status = ? WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, status, id)
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

// Additional helper methods

func (r *todoRepository) FindByUserIdAndStatus(ctx context.Context, userId uint, status string, page, limit int) ([]todo.Todo, int64, error) {
	offset := (page - 1) * limit

	var total int64
	countQuery := "SELECT COUNT(*) FROM todos WHERE user_id = ? AND status = ?"
	if err := r.db.QueryRowContext(ctx, countQuery, userId, status).Scan(&total); err != nil {
		return nil, 0, errors.Internal("Failed to count todos")
	}

	query := `
       SELECT id, user_id, title, description, status, due_date, created_at, updated_at
       FROM todos 
       WHERE user_id = ? AND status = ?
       ORDER BY created_at DESC
       LIMIT ? OFFSET ?
   `

	rows, err := r.db.QueryContext(ctx, query, userId, status, limit, offset)
	if err != nil {
		return nil, 0, errors.Internal("Failed to get todos")
	}
	defer rows.Close()

	var todos []todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Internal("Failed to scan todo")
		}
		todos = append(todos, t)
	}

	return todos, total, nil
}

func (r *todoRepository) FindUpcoming(ctx context.Context, userId uint, limit int) ([]todo.Todo, error) {
	query := `
       SELECT id, user_id, title, description, status, due_date, created_at, updated_at
       FROM todos 
       WHERE user_id = ? 
         AND status = 'pending'
         AND due_date > NOW()
       ORDER BY due_date ASC
       LIMIT ?
   `

	rows, err := r.db.QueryContext(ctx, query, userId, limit)
	if err != nil {
		return nil, errors.Internal("Failed to get upcoming todos")
	}
	defer rows.Close()

	var todos []todo.Todo
	for rows.Next() {
		var t todo.Todo
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, errors.Internal("Failed to scan todo")
		}
		todos = append(todos, t)
	}

	return todos, nil
}
