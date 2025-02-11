package user

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
