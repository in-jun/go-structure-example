package user

import (
	"net/http"

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
	users := r.Group("/users")
	{
		users.Use(middleware.Auth())
		users.GET("/me", h.GetMyInfo)
		users.PUT("/me", h.UpdateMyInfo)
		users.DELETE("/me", h.DeleteMyAccount)
	}
}

func (h *Handler) GetMyInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateMyInfo(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.Update(c.Request.Context(), userID, req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *Handler) DeleteMyAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
