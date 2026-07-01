package service

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
)

type StatisticStorage interface {
	GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error)
	GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error)
}

type StatisticService struct {
	storage StatisticStorage
}

func NewStatisticService(storage StatisticStorage) *StatisticService {
	return &StatisticService{storage: storage}
}

func (s *StatisticService) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	if email == "" {
		return nil, apperr.NewValidationError("email is required")
	}

	return s.storage.GetUserOrderHistory(ctx, email)
}

func (s *StatisticService) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	return s.storage.GetPopularProducts(ctx)
}
