package http

import (
	"net/http"
	"strconv"

	"github.com/in-jun/go-structure-example/internal/bid/application"
	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
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

	mux.Handle("GET /api/v1/auctions/{auction_id}/bids", mw(http.HandlerFunc(h.ListBids)))
	mux.Handle("GET /api/v1/auctions/{auction_id}/bids/highest", mw(http.HandlerFunc(h.GetHighest)))
	mux.Handle("POST /api/v1/auctions/{auction_id}/bids", mw(gatewayAuth(http.HandlerFunc(h.PlaceBid))))
}

func (h *Handler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	auctionID := r.PathValue("auction_id")

	var req PlaceBidRequest
	if err := server.Bind(r, &req); err != nil {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	userID := server.UserID(r)
	result, err := h.commands.PlaceBid(r.Context(), command.PlaceBid{
		UserID:    userID,
		AuctionID: auctionID,
		Amount:    req.Amount,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusCreated, toPlaceBidResponse(result))
}

func (h *Handler) ListBids(w http.ResponseWriter, r *http.Request) {
	auctionID := r.PathValue("auction_id")
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

	result, err := h.queries.ListBids(r.Context(), query.ListBids{
		AuctionID: auctionID,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toListResponse(result))
}

func (h *Handler) GetHighest(w http.ResponseWriter, r *http.Request) {
	auctionID := r.PathValue("auction_id")

	result, err := h.queries.GetHighest(r.Context(), query.GetHighest{
		AuctionID: auctionID,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	server.JSON(w, http.StatusOK, toGetResponse(result))
}
