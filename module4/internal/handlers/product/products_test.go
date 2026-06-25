package product

import (
	"bytes"
	"context"
	"errors"
	"go-pet-shop/internal/handlers/product/mocks"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mockedStorage := new(mocks.Products)
	mockedStorage.On("GetAllProducts", mock.Anything).Return([]models.Product{
		{ID: 1, Name: "Dog Food"},
	}, nil)

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mockedStorage)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mockedStorage := new(mocks.Products)
	mockedStorage.On("GetAllProducts", mock.Anything).Return(nil, errors.New("DB error"))

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	productToSave := models.Product{
		ID:    1,
		Name:  "testProduct",
		Price: 45.99,
		Stock: 120,
	}

	mockedStorage := new(mocks.Products)
	mockedStorage.On("CreateProduct", mock.Anything, productToSave).Return(1, nil)

	bodyJson := `{ "ID": 1, "Name": "testProduct", "Price": 45.99, "Stock": 120 }`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyJson))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	mockedStorage := new(mocks.Products)

	bodyJson := `{ invalid JSON }`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyJson))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	productToSave := models.Product{
		ID:    1,
		Name:  "testProduct",
		Price: 45.99,
		Stock: 120,
	}

	mockedStorage := new(mocks.Products)
	mockedStorage.On("CreateProduct", mock.Anything, productToSave).Return(0, errors.New("server error"))

	bodyJson := `{ "ID": 1, "Name": "testProduct", "Price": 45.99, "Stock": 120 }`

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyJson))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	productToSave := models.Product{
		ID:    1,
		Name:  "productToUpdate",
		Price: 45.99,
		Stock: 120,
	}

	mockedStorage := new(mocks.Products)
	mockedStorage.On("UpdateProduct", mock.Anything, productToSave).Return(nil)

	bodyJson := `{ "ID": 1, "Name": "productToUpdate", "Price": 45.99, "Stock": 120 }`

	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBufferString(bodyJson))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	mockedStorage := new(mocks.Products)

	bodyJson := `{ "ID": 1, "Name": "", "Price": 45.99, "Stock": 120 }`

	req := httptest.NewRequest(http.MethodPost, "/products/1", bytes.NewBufferString(bodyJson))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	productToSave := models.Product{
		ID:    1,
		Name:  "productToUpdate",
		Price: 45.99,
		Stock: 120,
	}

	mockedStorage := new(mocks.Products)
	mockedStorage.On("UpdateProduct", mock.Anything, productToSave).Return(errors.New("server error"))

	bodyJson := `{ "ID": 1, "Name": "productToUpdate", "Price": 45.99, "Stock": 120 }`

	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBufferString(bodyJson))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	mockedStorage := new(mocks.Products)
	mockedStorage.On("DeleteProduct", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	mockedStorage := new(mocks.Products)

	req := httptest.NewRequest(http.MethodDelete, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	mockedStorage := new(mocks.Products)
	mockedStorage.On("DeleteProduct", mock.Anything, 1).Return(errors.New("server error"))

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mockedStorage)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
