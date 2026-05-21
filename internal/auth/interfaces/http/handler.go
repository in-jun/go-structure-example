package http

import (
	"net/http"

	"github.com/in-jun/go-structure-example/internal/auth/application"
	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type Handler struct {
	commands      application.CommandUseCase
	queries       application.QueryUseCase
	validateToken middleware.TokenValidator
}

func NewHandler(commands application.CommandUseCase, queries application.QueryUseCase, validateToken middleware.TokenValidator) *Handler {
	return &Handler{commands: commands, queries: queries, validateToken: validateToken}
}

func (h *Handler) RegisterRoutes(mux *server.Router, mw server.Middleware) {
	authMw := middleware.Auth(h.validateToken)

	mux.Handle("POST /api/v1/auth/register", mw(http.HandlerFunc(h.Register)))
	mux.Handle("POST /api/v1/auth/login", mw(http.HandlerFunc(h.Login)))
	mux.Handle("POST /api/v1/auth/refresh", mw(http.HandlerFunc(h.Refresh)))
	mux.Handle("POST /api/v1/auth/logout", mw(authMw(http.HandlerFunc(h.Logout))))
	mux.Handle("POST /api/v1/auth/logout/all", mw(authMw(http.HandlerFunc(h.LogoutAll))))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	if err := h.commands.Register(r.Context(), command.Register{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, MessageResponse{Message: "Registration successful"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	result, err := h.commands.Login(r.Context(), command.Login{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toLoginResponse(result))
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	result, err := h.commands.Refresh(r.Context(), command.Refresh{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toRefreshResponse(result))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	jti := server.TokenJTI(r)

	if err := h.commands.Logout(r.Context(), command.Logout{
		RefreshToken:   req.RefreshToken,
		AccessTokenJTI: jti,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, MessageResponse{Message: "Logout successful"})
}

func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	userID := server.UserID(r)

	if err := h.commands.LogoutAll(r.Context(), command.LogoutAll{
		UserID: userID,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, MessageResponse{Message: "All sessions logged out successfully"})
}
