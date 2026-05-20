package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/todo/application"
	"github.com/in-jun/go-structure-example/internal/todo/application/command"
	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type Handler struct {
	commands      application.CommandUseCase
	queries       application.QueryUseCase
	validateToken middleware.TokenValidator
}

func NewHandler(commands application.CommandUseCase, queries application.QueryUseCase, validateToken middleware.TokenValidator) *Handler {
	return &Handler{commands: commands, queries: queries, validateToken: validateToken}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	todos := r.Group("/todos")
	todos.Use(middleware.Auth(h.validateToken))
	{
		todos.GET("", h.GetList)
		todos.POST("", h.Create)
		todos.GET("/:id", h.Get)
		todos.PUT("/:id", h.Update)
		todos.PATCH("/:id/status", h.UpdateStatus)
		todos.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.commands.Create(c.Request.Context(), command.Create{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.ID})
}

func (h *Handler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID := c.GetUint("user_id")

	result, err := h.queries.ListTodos(c.Request.Context(), query.ListTodos{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toTodoListResponse(result))
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("Invalid todo ID"))
		return
	}

	userID := c.GetUint("user_id")
	result, err := h.queries.GetTodo(c.Request.Context(), query.GetTodo{
		UserID: userID,
		TodoID: uint(id),
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toTodoResponse(result))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("Invalid todo ID"))
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.commands.Update(c.Request.Context(), command.Update{
		UserID:      userID,
		TodoID:      uint(id),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("Invalid todo ID"))
		return
	}

	var req UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.commands.UpdateStatus(c.Request.Context(), command.UpdateStatus{
		UserID: userID,
		TodoID: uint(id),
		Status: entity.Status(req.Status),
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("Invalid todo ID"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.commands.Delete(c.Request.Context(), command.Delete{
		UserID: userID,
		TodoID: uint(id),
	}); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
