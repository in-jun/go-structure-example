package http

import (
	"net/http"
	"strconv"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
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

func (h *Handler) RegisterRoutes(mux *server.Router, mw server.Middleware) {
	authMw := middleware.Auth(h.validateToken)

	mux.Handle("GET /api/v1/todos", mw(authMw(http.HandlerFunc(h.GetList))))
	mux.Handle("POST /api/v1/todos", mw(authMw(http.HandlerFunc(h.Create))))
	mux.Handle("GET /api/v1/todos/{id}", mw(authMw(http.HandlerFunc(h.Get))))
	mux.Handle("PUT /api/v1/todos/{id}", mw(authMw(http.HandlerFunc(h.Update))))
	mux.Handle("PATCH /api/v1/todos/{id}/status", mw(authMw(http.HandlerFunc(h.UpdateStatus))))
	mux.Handle("DELETE /api/v1/todos/{id}", mw(authMw(http.HandlerFunc(h.Delete))))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	result, err := h.commands.Create(r.Context(), command.Create{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, map[string]string{"id": result.ID})
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(server.QueryDefault(r, "page", "1"))
	limit, _ := strconv.Atoi(server.QueryDefault(r, "limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}
	userID := server.UserID(r)

	result, err := h.queries.GetList(r.Context(), query.List{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toTodoListResponse(result))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := server.PathParam(r, "id")
	userID := server.UserID(r)
	result, err := h.queries.Get(r.Context(), query.Get{
		UserID: userID,
		TodoID: id,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toTodoResponse(result))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := server.PathParam(r, "id")

	var req UpdateTodoRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	if err := h.commands.Update(r.Context(), command.Update{
		UserID:      userID,
		TodoID:      id,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := server.PathParam(r, "id")

	var req UpdateTodoStatusRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	if err := h.commands.UpdateStatus(r.Context(), command.UpdateStatus{
		UserID: userID,
		TodoID: id,
		Status: entity.Status(req.Status),
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := server.PathParam(r, "id")

	userID := server.UserID(r)
	if err := h.commands.Delete(r.Context(), command.Delete{
		UserID: userID,
		TodoID: id,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
