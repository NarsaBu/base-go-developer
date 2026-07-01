package service

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
)

type OrderStorage interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

type OrderService struct {
	storage OrderStorage
}

func NewOrderService(storage OrderStorage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	if order.UserID <= 0 {
		return 0, apperr.NewValidationError("user id is required and must be positive")
	}

	return s.storage.CreateOrder(ctx, order)
}

func (s *OrderService) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	if orderItem.OrderID <= 0 {
		return apperr.NewValidationError("order id must be positive")
	}
	if orderItem.ProductID <= 0 {
		return apperr.NewValidationError("product id is required and must be positive")
	}
	if orderItem.Quantity <= 0 {
		return apperr.NewValidationError("quantity should be a positive natural number")
	}

	return s.storage.AddOrderItem(ctx, orderItem)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	if id <= 0 {
		return models.Order{}, apperr.NewValidationError("order id must be positive")
	}
	return s.storage.GetOrderByID(ctx, id)
}

func (s *OrderService) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	if email == "" {
		return nil, apperr.NewValidationError("email is required")
	}
	return s.storage.GetOrdersByUserEmail(ctx, email)
}

func (s *OrderService) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	if orderID <= 0 {
		return nil, apperr.NewValidationError("order id must be positive")
	}
	return s.storage.GetOrderItemsByOrderID(ctx, orderID)
}
