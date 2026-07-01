package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/apperr"
	"go-pet-shop/internal/models"

	"github.com/lib/pq"
)

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int, error) {
	const fn = "storage.postgres.user.CreateUser"

	var id int
	err := s.db.QueryRow(ctx, `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`, user.Name, user.Email).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, apperr.ErrAliasAlreadyExists
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const fn = "storage.postgres.user.GetUserByEmail"

	var user models.User

	err := s.db.QueryRow(ctx, `SELECT id, name, email FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (s *Storage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	const fn = "storage.postgres.user.GetAllUsers"

	rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return users, nil
}
