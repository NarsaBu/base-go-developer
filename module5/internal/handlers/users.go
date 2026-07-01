package handlers

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Users interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type UserHandler struct {
	log     *slog.Logger
	service Users
}

func NewUserHandler(log *slog.Logger, service Users) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.CreateUser"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))
	log.Info("Creating new user", slog.String("url", r.URL.String()))

	var user models.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad request", "message": "Invalid JSON payload"})
		return
	}

	userID, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		if apperr.IsValidationError(err) {
			log.Warn("validation failed", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}

		if errors.Is(err, apperr.ErrEmailAlreadyExists) {
			log.Warn("attempt to create user with existing email", slog.String("email", user.Email))
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, map[string]string{"error": "Conflict", "message": "User with this email already exists"})
			return
		}

		log.Error("failed to create user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to create user"})
		return
	}

	log.Info("User created successfully", slog.Int("id", userID), slog.String("name", user.Name))
	user.ID = userID
	render.JSON(w, r, map[string]interface{}{"status": "User created successfully", "id": userID, "user": user})
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.GetUserByEmail"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

	email := chi.URLParam(r, "email")

	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		if apperr.IsValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Bad request", "message": err.Error()})
			return
		}
		if errors.Is(err, apperr.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{"error": "Not found", "message": fmt.Sprintf("User with email %s does not exist", email)})
			return
		}
		log.Error("failed to get user by email", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve user"})
		return
	}

	render.JSON(w, r, user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.users.GetAllUsers"
	log := h.log.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		log.Error("failed to get users", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Internal server error", "message": "Failed to retrieve users"})
		return
	}

	render.JSON(w, r, users)
}
