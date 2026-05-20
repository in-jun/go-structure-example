package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/auth/application"
	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
)

type Handler struct {
	commands      application.CommandUseCase
	validateToken middleware.TokenValidator
}

func NewHandler(commands application.CommandUseCase, validateToken middleware.TokenValidator) *Handler {
	return &Handler{commands: commands, validateToken: validateToken}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", middleware.Auth(h.validateToken), h.Logout)
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	if err := h.commands.Register(c.Request.Context(), command.Register{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{Message: "Registration successful"})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	result, err := h.commands.Login(c.Request.Context(), command.Login{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toAuthResponse(result))
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	result, err := h.commands.Refresh(c.Request.Context(), command.Refresh{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toRefreshResponse(result))
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	jti := c.GetString("jti")

	if err := h.commands.Logout(c.Request.Context(), command.Logout{
		RefreshToken:   req.RefreshToken,
		AccessTokenJTI: jti,
	}); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Logout successful"})
}
