package todo

import "time"

type CreateTodoRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTodoStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed"`
}

type TodoResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
	Total int64          `json:"total"`
}
