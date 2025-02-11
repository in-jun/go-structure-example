package auth

import (
	"net/http"
	"strings"

	"github.com/in-jun/go-structure-example/internal/pkg/middleware"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", middleware.Auth(), h.Logout)
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	if err := h.service.Register(c.Request.Context(), req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Registration successful"})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	res, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Error(errors.Unauthorized("Missing authorization header"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.Error(errors.Unauthorized("Invalid authorization header"))
		return
	}

	tokenString := parts[1]
	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		c.Error(errors.Unauthorized("Invalid token"))
		return
	}

	if err := h.service.Logout(c.Request.Context(), claims.UserID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Logout successful"})
}
