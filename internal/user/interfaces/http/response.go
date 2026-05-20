package http

import (
	"time"

	"github.com/in-jun/go-structure-example/internal/user/application/query"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func toUserResponse(r *query.Result) *UserResponse {
	return &UserResponse{
		ID:        r.ID,
		Email:     r.Email,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}
