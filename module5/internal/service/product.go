package service

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
)

type ProductStorage interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, id int) (models.Product, error)
}

type ProductService struct {
	storage ProductStorage
}

func NewProductService(storage ProductStorage) *ProductService {
	return &ProductService{storage: storage}
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.storage.GetAllProducts(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	if product.Name == "" {
		return 0, apperr.NewValidationError("product name is required")
	}
	if product.Price < 0 {
		return 0, apperr.NewValidationError("product price cannot be negative")
	}
	if product.Stock < 0 {
		return 0, apperr.NewValidationError("product stock cannot be negative")
	}

	return s.storage.CreateProduct(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return apperr.NewValidationError("invalid product ID")
	}
	return s.storage.DeleteProduct(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product models.Product) error {
	if product.ID <= 0 {
		return apperr.NewValidationError("invalid product ID")
	}
	if product.Name == "" {
		return apperr.NewValidationError("product name is required")
	}
	if product.Price < 0 {
		return apperr.NewValidationError("product price cannot be negative")
	}
	if product.Stock < 0 {
		return apperr.NewValidationError("product stock cannot be negative")
	}

	return s.storage.UpdateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	if id <= 0 {
		return models.Product{}, apperr.NewValidationError("invalid product ID")
	}
	return s.storage.GetProductByID(ctx, id)
}
