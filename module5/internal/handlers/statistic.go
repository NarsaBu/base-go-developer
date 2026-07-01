package handlers

import (
	"context"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Statistic interface {
	GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error)
	GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error)
}

type StatisticHandler struct {
	log     *slog.Logger
	storage Statistic
}

func NewStatisticHandler(log *slog.Logger, storage Statistic) *StatisticHandler {
	return &StatisticHandler{
		log:     log,
		storage: storage,
	}
}

func (h *StatisticHandler) GetUserOrderHistory(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.statisitc.GetUserOrderHistory"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting user order history", slog.String("url", r.URL.String()))

	email := r.URL.Query().Get("email")
	if email == "" {
		log.Error("empty id in URL")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email is required",
		})
		return
	}

	orderHistory, err := h.storage.GetUserOrderHistory(r.Context(), email)

	if err != nil {
		log.Error("failed to get user order history", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve user order history",
		})
		return
	}

	log.Info("Retrieved user order history successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, orderHistory)
}

func (h *StatisticHandler) GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.statisitc.GetPopularProducts"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Getting popular products", slog.String("url", r.URL.String()))

	popularProducts, err := h.storage.GetPopularProducts(r.Context())

	if err != nil {
		log.Error("failed to get popular products", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve popular products",
		})
		return
	}

	log.Info("Retrieved popular products successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, popularProducts)
}
