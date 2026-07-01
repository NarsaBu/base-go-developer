package handlers

import (
	"context"
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
	storage Orders
}

func NewOrderHandler(log *slog.Logger, storage Orders) *OrderHandler {
	return &OrderHandler{
		log:     log,
		storage: storage,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.CreateOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new order", slog.String("url", r.URL.String()))

	var order models.Order
	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if order.UserID <= 0 {
		log.Error("user id is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User id is required",
		})
		return
	}

	orderId, err := h.storage.CreateOrder(r.Context(), order)
	if err != nil {
		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info("Order created successfully",
		slog.Int("id", orderId),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order created successfully",
		"id":     orderId,
	})
}

func (h *OrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.AddOrderItem"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Adding new order item", slog.String("url", r.URL.String()))

	var orderItem models.OrderItem
	if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	orderIdStr := chi.URLParam(r, "id")
	if orderIdStr == "" {
		log.Error("empty id in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Id is required",
		})
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", orderIdStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a number",
		})
		return
	}

	orderItem.OrderID = orderId

	if orderItem.ProductID <= 0 {
		log.Error("product id is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "product id is required",
		})
		return
	}

	if orderItem.Quantity <= 0 {
		log.Error("quantity is not valid")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "quantity should contains positive natural number",
		})
		return
	}

	err = h.storage.AddOrderItem(r.Context(), orderItem)
	if err != nil {
		log.Error("failed to create order item", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order item",
		})
		return
	}

	log.Info("Order item created successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "Order created successfully",
	})
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrderByID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting order by id", slog.String("url", r.URL.String()))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty id in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a number",
		})
		return
	}

	order, err := h.storage.GetOrderByID(r.Context(), id)

	if err != nil {
		log.Error("failed to get order by id", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order",
		})
		return
	}

	log.Info("Retrieved order successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, order)
}

func (h *OrderHandler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrdersByUserEmail"
	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.Query().Get("email")
	if email == "" {
		log.Error("empty email in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email is required",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)

	if err != nil {
		log.Error("failed to get orders", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders",
		})
		return
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()),
		slog.Int("count", len(orders)),
	)

	render.JSON(w, r, orders)
}

func (h *OrderHandler) GetOrderItemsByOrderID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.orders.GetOrderItemsByOrderID"
	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("empty id in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a number",
		})
		return
	}

	orderItems, err := h.storage.GetOrderItemsByOrderID(r.Context(), id)

	if err != nil {
		log.Error("failed to get order items", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve order items",
		})
		return
	}

	log.Info("Retrieved order items successfully",
		slog.String("url", r.URL.String()),
		slog.Int("count", len(orderItems)),
	)

	render.JSON(w, r, orderItems)
}
