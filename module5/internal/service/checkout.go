package service

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
)

type CheckoutStorage interface {
	PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (int, error)
}

type CheckoutService struct {
	storage CheckoutStorage
}

func NewCheckoutService(storage CheckoutStorage) *CheckoutService {
	return &CheckoutService{storage: storage}
}

func (s *CheckoutService) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
	if userEmail == "" {
		return 0, apperr.NewValidationError("user email is required")
	}

	if len(items) == 0 {
		return 0, apperr.NewValidationError("order must contain at least one item")
	}

	for _, item := range items {
		if item.ProductID <= 0 {
			return 0, apperr.NewValidationError(
				"product id is required and must be positive",
			)
		}
		if item.Quantity <= 0 {
			return 0, apperr.NewValidationError(
				"quantity must be a positive natural number",
			)
		}
	}

	return s.storage.PlaceOrder(ctx, userEmail, items)
}
