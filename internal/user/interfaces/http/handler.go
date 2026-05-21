package http

import (
	"net/http"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
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

func (h *Handler) RegisterRoutes(mux *server.Router, mw server.Middleware) {
	authMw := middleware.Auth(h.validateToken)

	mux.Handle("GET /api/v1/users/me", mw(authMw(http.HandlerFunc(h.GetMe))))
	mux.Handle("PATCH /api/v1/users/me/profile", mw(authMw(http.HandlerFunc(h.UpdateProfile))))
	mux.Handle("PATCH /api/v1/users/me/password", mw(authMw(http.HandlerFunc(h.UpdatePassword))))
	mux.Handle("DELETE /api/v1/users/me", mw(authMw(http.HandlerFunc(h.DeleteMe))))
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := server.UserID(r)

	result, err := h.queries.GetProfile(r.Context(), query.Get{UserID: userID})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toUserResponse(result))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	if err := h.commands.UpdateProfile(r.Context(), command.UpdateProfile{
		UserID: userID,
		Name:   req.Name,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, MessageResponse{Message: "Profile updated successfully"})
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	if err := h.commands.UpdatePassword(r.Context(), command.UpdatePassword{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, MessageResponse{Message: "Password updated successfully"})
}

func (h *Handler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID := server.UserID(r)
	if err := h.commands.Delete(r.Context(), command.Delete{UserID: userID}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, MessageResponse{Message: "Account deleted successfully"})
}
