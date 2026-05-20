package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/user/application"
	"github.com/in-jun/go-structure-example/internal/user/application/command"
	"github.com/in-jun/go-structure-example/internal/user/application/query"
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
	users := r.Group("/users")
	users.Use(middleware.Auth(h.validateToken))
	{
		users.GET("/me", h.GetMe)
		users.PATCH("/me/profile", h.UpdateProfile)
		users.PATCH("/me/password", h.UpdatePassword)
		users.DELETE("/me", h.DeleteMe)
	}
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetUint("user_id")

	result, err := h.queries.GetProfile(c.Request.Context(), query.Get{UserID: userID})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toUserResponse(result))
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.commands.UpdateProfile(c.Request.Context(), command.UpdateProfile{
		UserID: userID,
		Name:   req.Name,
	}); err != nil {
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
	if err := h.commands.UpdatePassword(c.Request.Context(), command.UpdatePassword{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Password updated successfully"})
}

func (h *Handler) DeleteMe(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.commands.Delete(c.Request.Context(), command.Delete{UserID: userID}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Account deleted successfully"})
}
