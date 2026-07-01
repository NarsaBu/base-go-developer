package handlers

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

type OrderHandler struct {
	log     *slog.Logger
	service Orders
}

func NewOrderHandler(log *slog.Logger, service Orders) *OrderHandler {
	return &OrderHandler{
		log:     log,
		service: service,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.CreateOrder"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))
	log.Info("Creating new order", slog.String("url", r.URL.String()))

	var order models.Order
	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Invalid JSON payload"})
		return
	}

	orderID, err := h.service.CreateOrder(r.Context(), order)
	if err != nil {
		if apperr.IsValidationError(err) {
			log.Warn("validation failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to create order"})
		return
	}

	log.Info("Order created successfully", slog.Int("id", orderID), slog.String("url", r.URL.String()))
	render.JSON(w, r, map[string]interface{}{"status": "Order created successfully", "id": orderID})
}

func (h *OrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.AddOrderItem"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))
	log.Info("Adding new order item", slog.String("url", r.URL.String()))

	var orderItem models.OrderItem
	if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Invalid JSON payload"})
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Id is required"})
		return
	}

	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Order ID must be a number"})
		return
	}
	orderItem.OrderID = orderID

	if err := h.service.AddOrderItem(r.Context(), orderItem); err != nil {
		if apperr.IsValidationError(err) {
			log.Warn("validation failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		log.Error("failed to create order item", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to create order item"})
		return
	}

	log.Info("Order item created successfully", slog.String("url", r.URL.String()))
	render.JSON(w, r, map[string]interface{}{"status": "Order item created successfully"})
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrderByID"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))
	log.Info("Getting order by id", slog.String("url", r.URL.String()))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Order ID must be a number"})
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), id)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		log.Error("failed to get order by id", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve order"})
		return
	}

	log.Info("Retrieved order successfully", slog.String("url", r.URL.String()))
	render.JSON(w, r, order)
}

func (h *OrderHandler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrdersByUserEmail"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

	email := r.URL.Query().Get("email")

	orders, err := h.service.GetOrdersByUserEmail(r.Context(), email)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		log.Error("failed to get orders", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve orders"})
		return
	}

	log.Info("Retrieved orders successfully", slog.String("url", r.URL.String()), slog.Int("count", len(orders)))
	render.JSON(w, r, orders)
}

func (h *OrderHandler) GetOrderItemsByOrderID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrderItemsByOrderID"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Order ID must be a number"})
		return
	}

	orderItems, err := h.service.GetOrderItemsByOrderID(r.Context(), id)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		log.Error("failed to get order items", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve order items"})
		return
	}

	log.Info("Retrieved order items successfully", slog.String("url", r.URL.String()), slog.Int("count", len(orderItems)))
	render.JSON(w, r, orderItems)
}
