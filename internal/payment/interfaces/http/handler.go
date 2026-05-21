package http

import (
	"io"
	"net/http"

	"github.com/in-jun/go-structure-example/internal/payment/application"
	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	"github.com/in-jun/go-structure-example/internal/payment/application/query"
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

	mux.Handle("GET /api/v1/payments/{id}", mw(gatewayAuth(http.HandlerFunc(h.GetPayment))))
	mux.Handle("POST /api/v1/payments/{id}/confirm", mw(gatewayAuth(http.HandlerFunc(h.ConfirmPayment))))
	mux.Handle("POST /api/v1/payments/{id}/refund", mw(gatewayAuth(http.HandlerFunc(h.RefundPayment))))
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	result, err := h.queries.GetPayment(r.Context(), query.GetPayment{
		PaymentID: id,
	})
	if err != nil {
		middleware.HandleError(w, err)
		return
	}

	if result.WinnerID != userID {
		middleware.HandleError(w, errors.Forbidden("Not authorized"))
		return
	}

	server.JSON(w, http.StatusOK, toGetResponse(result))
}

func (h *Handler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	if err := h.commands.ConfirmPayment(r.Context(), command.ConfirmPayment{
		UserID:    userID,
		PaymentID: id,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := server.UserID(r)

	var req RefundRequest
	if err := server.Bind(r, &req); err != nil && err != io.EOF {
		middleware.HandleError(w, errors.BadRequest("Invalid request format"))
		return
	}

	if err := h.commands.RefundPayment(r.Context(), command.RefundPayment{
		UserID:    userID,
		PaymentID: id,
		Reason:    req.Reason,
	}); err != nil {
		middleware.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
