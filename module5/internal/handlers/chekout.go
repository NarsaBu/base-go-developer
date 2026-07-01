package handlers

import (
	"context"
	"errors"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Checkout interface {
	PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error)
}

type CheckoutRequest struct {
	UserEmail string             `json:"userEmail"`
	Items     []models.OrderItem `json:"items"`
}

type CheckoutHandler struct {
	log     *slog.Logger
	service Checkout
}

func NewCheckoutHandler(log *slog.Logger, service Checkout) *CheckoutHandler {
	return &CheckoutHandler{
		log:     log,
		service: service,
	}
}

func (h *CheckoutHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.checkout.PlaceOrder"
	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	log.Info("Placing new order", slog.String("url", r.URL.String()))

	var requestBody CheckoutRequest
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	orderID, err := h.service.PlaceOrder(r.Context(), requestBody.UserEmail, requestBody.Items)
	if err != nil {
		if apperr.IsValidationError(err) {
			log.Warn("validation failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, apperr.ErrUserNotFound) {
			log.Warn("user not found for checkout", slog.String("email", requestBody.UserEmail))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, apperr.ErrInsufficientStock) {
			log.Warn("insufficient stock during checkout", slog.Any("error", err))
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, map[string]string{
				"error":   "Conflict",
				"message": err.Error(),
			})
			return
		}

		log.Error("failed to place order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to place order",
		})
		return
	}

	log.Info("Order placed successfully",
		slog.Int("id", orderID),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order placed successfully",
		"id":     orderID,
	})
}
