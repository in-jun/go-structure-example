package mysql

import (
	"context"
	"database/sql"

	"github.com/in-jun/go-structure-example/internal/app/user"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *user.User) error {
	query := `
       INSERT INTO users (email, password, name)
       VALUES (?, ?, ?)
   `

	result, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return errors.Internal("Failed to create user")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.Internal("Failed to get user ID")
	}

	user.ID = uint(id)
	return nil
}

func (r *userRepository) FindById(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	query := `
       SELECT id, email, password, name, created_at, updated_at
       FROM users WHERE id = ?
   `

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get user")
	}

	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `
       SELECT id, email, password, name, created_at, updated_at
       FROM users WHERE email = ?
   `

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get user by email")
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	query := `
       UPDATE users 
       SET name = ?, password = ?
       WHERE id = ?
   `

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Password, user.ID)
	if err != nil {
		return errors.Internal("Failed to update user")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("User not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Internal("Failed to delete user")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}

	if rows == 0 {
		return errors.NotFound("User not found")
	}

	return nil
}
