package todo

import (
	"net/http"
	"strconv"

	"github.com/in-jun/go-structure-example/internal/pkg/middleware"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	todos := r.Group("/todos")
	{
		todos.Use(middleware.Auth())
		todos.GET("", h.GetList)
		todos.POST("", h.Create)
		todos.GET("/:id", h.GetTodoDetail)
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
	res, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	userID := c.GetUint("user_id")
	res, err := h.service.GetList(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
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
	if err := h.service.Update(c.Request.Context(), userID, uint(id), req); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetTodoDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequest("Invalid todo ID"))
		return
	}

	userID := c.GetUint("user_id")
	todo, err := h.service.GetByID(c.Request.Context(), userID, uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, todo)
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
	if err := h.service.UpdateStatus(c.Request.Context(), userID, uint(id), req); err != nil {
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
	if err := h.service.Delete(c.Request.Context(), userID, uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
