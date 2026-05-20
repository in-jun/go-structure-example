package pg

import (
	"context"
	"database/sql"
	stderrors "errors"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db func(ctx context.Context) transaction.DBTX
}

func NewUserRepository(db func(ctx context.Context) transaction.DBTX) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (id, email, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db(ctx).ExecContext(ctx, query,
		user.ID(), user.Email(), user.HashedPassword(), user.Name(), user.CreatedAt(), user.UpdatedAt())
	if err != nil {
		return errors.Internal("Failed to create user")
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := "SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1"
	var id, e, password, name string
	var createdAt, updatedAt time.Time
	err := r.db(ctx).QueryRowContext(ctx, query, email).Scan(&id, &e, &password, &name, &createdAt, &updatedAt)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get user by email")
	}
	u, err := entity.ReconstructUser(id, e, password, name, createdAt, updatedAt)
	if err != nil {
		return nil, errors.Internal("Failed to reconstruct user")
	}
	return u, nil
}
