package handlers

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Products interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int) (models.Product, error)
}

type ProductHandler struct {
	log     *slog.Logger
	service Products // Раньше было storage
}

func NewProductHandler(log *slog.Logger, service Products) *ProductHandler {
	return &ProductHandler{log: log, service: service}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve products"})
		return
	}
	render.JSON(w, r, items)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := render.DecodeJSON(r.Body, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Invalid JSON payload"})
		return
	}

	productID, err := h.service.CreateProduct(r.Context(), product)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to create product"})
		return
	}

	product.ID = productID
	render.JSON(w, r, map[string]interface{}{"status": "Product created successfully", "id": productID, "product": product})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr) // Конвертация HTTP-параметра остается в хендлере
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Product ID must be a number"})
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}
		if errors.Is(err, apperr.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{"error": "Not found", "message": fmt.Sprintf("Product with ID %d does not exist", id)})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to delete product"})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "Product deleted successfully", "id": id})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Product ID must be a number"})
		return
	}

	var product models.Product
	if err := render.DecodeJSON(r.Body, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Invalid JSON payload"})
		return
	}
	product.ID = id

	if err := h.service.UpdateProduct(r.Context(), product); err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}
		if errors.Is(err, apperr.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{"error": "Not found", "message": fmt.Sprintf("Product with ID %d does not exist", id)})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to update product"})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "Product updated successfully", "id": id, "product": product})
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Product ID must be a number"})
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}
		if errors.Is(err, apperr.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{"error": "Not found", "message": fmt.Sprintf("Product with ID %d does not exist", id)})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve product"})
		return
	}

	render.JSON(w, r, product)
}
