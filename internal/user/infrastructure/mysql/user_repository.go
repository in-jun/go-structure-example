package mysql

import (
	"context"
	"database/sql"
	stderrors "errors"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/in-jun/go-structure-example/internal/user/domain"
	"github.com/in-jun/go-structure-example/internal/user/domain/entity"
)

var _ domain.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db func(ctx context.Context) transaction.DBTX
}

func NewUserRepository(db func(ctx context.Context) transaction.DBTX) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	query := "SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = ?"
	var uid uint
	var email, password, name string
	var createdAt, updatedAt time.Time
	err := r.db(ctx).QueryRowContext(ctx, query, id).Scan(&uid, &email, &password, &name, &createdAt, &updatedAt)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get user")
	}
	u, err := entity.ReconstructUser(uid, email, password, name, createdAt, updatedAt)
	if err != nil {
		return nil, errors.Internal("Failed to reconstruct user")
	}
	return u, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := "UPDATE users SET name = ?, password = ? WHERE id = ?"
	result, err := r.db(ctx).ExecContext(ctx, query, user.Name(), user.HashedPassword(), user.ID())
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
	result, err := r.db(ctx).ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
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
