package user

import (
	"net/http"

	"github.com/in-jun/go-structure-example/internal/pkg/middleware"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.Auth())
	{
		users.GET("/me", h.GetMe)
		users.PATCH("/me/profile", h.UpdateProfile)
		users.PATCH("/me/password", h.UpdatePassword)
		users.DELETE("/me", h.DeleteMe)
	}
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	u, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, u)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.UpdateProfile(c.Request.Context(), userID, req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Profile updated successfully"})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.UpdatePassword(c.Request.Context(), userID, req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Password updated successfully"})
}

func (h *Handler) DeleteMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), userID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Account deleted successfully"})
}
