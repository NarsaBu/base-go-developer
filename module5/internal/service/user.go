package service

import (
	"context"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type UserService struct {
	storage UserStorage
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (int, error) {
	if user.Name == "" {
		return 0, apperr.NewValidationError("user name is required")
	}
	if user.Email == "" {
		return 0, apperr.NewValidationError("user email is required")
	}

	return s.storage.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if email == "" {
		return models.User{}, apperr.NewValidationError("email is required")
	}
	return s.storage.GetUserByEmail(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.storage.GetAllUsers(ctx)
}
