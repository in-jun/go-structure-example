package http

import (
	"net/http"
	"strconv"

	"github.com/in-jun/go-structure-example/internal/auction/application"
	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type Handler struct {
	commands application.CommandUseCase
	queries  application.QueryUseCase
}

func NewHandler(commands application.CommandUseCase, queries application.QueryUseCase) *Handler {
	return &Handler{commands: commands, queries: queries}
}

func (h *Handler) RegisterRoutes(mux *server.Router, mw server.Middleware) {
	gatewayAuth := middleware.GatewayAuth()

	mux.Handle("GET /api/v1/auctions", mw(http.HandlerFunc(h.GetList)))
	mux.Handle("GET /api/v1/auctions/{id}", mw(http.HandlerFunc(h.GetByID)))
	mux.Handle("GET /api/v1/auctions/{id}/events", mw(http.HandlerFunc(h.GetEvents)))
	mux.Handle("POST /api/v1/auctions", mw(gatewayAuth(http.HandlerFunc(h.Create))))
	mux.Handle("POST /api/v1/auctions/{id}/open", mw(gatewayAuth(http.HandlerFunc(h.Open))))
	mux.Handle("POST /api/v1/auctions/{id}/close", mw(gatewayAuth(http.HandlerFunc(h.Close))))
	mux.Handle("POST /api/v1/auctions/{id}/cancel", mw(gatewayAuth(http.HandlerFunc(h.Cancel))))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	result, err := h.commands.Create(r.Context(), command.Create{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		StartPrice:  req.StartPrice,
		EndTime:     req.EndTime,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, toCreateResponse(result))
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(server.QueryDefault(r, "page", "1"))
	limit, _ := strconv.Atoi(server.QueryDefault(r, "limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	} else if limit > 100 {
		limit = 100
	}

	result, err := h.queries.GetList(r.Context(), query.List{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toListResponse(result))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.queries.GetByID(r.Context(), query.Get{
		AuctionID: id,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toGetResponse(result))
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.queries.GetEvents(r.Context(), query.EventHistory{
		AuctionID: id,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toEventHistoryResponse(result))
}

func (h *Handler) Open(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	if err := h.commands.Open(r.Context(), command.Open{
		UserID:    userID,
		AuctionID: id,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	if err := h.commands.Close(r.Context(), command.Close{
		UserID:    userID,
		AuctionID: id,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	if err := h.commands.Cancel(r.Context(), command.Cancel{
		UserID:    userID,
		AuctionID: id,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
