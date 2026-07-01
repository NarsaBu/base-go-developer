package handlers

import (
	"context"
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
	storage Checkout
}

func NewCheckoutHandler(log *slog.Logger, storage Checkout) *CheckoutHandler {
	return &CheckoutHandler{
		log:     log,
		storage: storage,
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

	if requestBody.UserEmail == "" {
		log.Error("user email is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User email is required",
		})
		return
	}

	for _, item := range requestBody.Items {
		if item.ProductID <= 0 {
			log.Error("order item product id is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": "Order item product id is required",
			})
			return
		}

		if item.Quantity <= 0 {
			log.Error("order item quantity is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": "Order item quantity is required",
			})
			return
		}
	}

	orderId, err := h.storage.PlaceOrder(r.Context(), requestBody.UserEmail, requestBody.Items)
	if err != nil {
		log.Error("failed to place order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to place order",
		})
		return
	}

	log.Info("Order placed successfully",
		slog.Int("id", orderId),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order placed successfully",
		"id":     orderId,
	})
}
