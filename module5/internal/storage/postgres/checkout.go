package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-pet-shop/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (orderID int, err error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("error while beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	userID, err := s.getUserIDByEmail(ctx, tx, userEmail)
	if err != nil {
		return 0, err
	}

	totalPrice, err := s.decrementStocks(ctx, tx, items)
	if err != nil {
		return 0, err
	}

	orderID, err = s.createOrder(ctx, tx, userID, totalPrice)
	if err != nil {
		return 0, err
	}

	if err = s.insertOrderItems(ctx, tx, orderID, items); err != nil {
		return 0, err
	}

	if err = s.insertTransaction(ctx, tx, orderID, totalPrice); err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("error while commiting transaction: %w", err)
	}

	return orderID, nil
}

func (s *Storage) getUserIDByEmail(ctx context.Context, tx pgx.Tx, email string) (int, error) {
	var userID int

	err := tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("user with email %s not found", email)
		}
		return 0, fmt.Errorf("error while getting user: %w", err)
	}

	return userID, nil
}

func (s *Storage) decrementStocks(ctx context.Context, tx pgx.Tx, items []models.OrderItem) (float64, error) {
	var totalPrice float64

	query := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1 RETURNING price`

	for _, item := range items {
		var price float64

		err := tx.QueryRow(ctx, query, item.Quantity, item.ProductID).Scan(&price)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return 0, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
			}
			return 0, fmt.Errorf("error while updating stock for product %d: %w", item.ProductID, err)
		}

		totalPrice += price * float64(item.Quantity)
	}

	return totalPrice, nil
}

func (s *Storage) createOrder(ctx context.Context, tx pgx.Tx, userID int, totalPrice float64) (int, error) {
	var orderID int

	query := `INSERT INTO orders (user_id, total_price, created_at) VALUES ($1, $2, $3) RETURNING id`

	err := tx.QueryRow(ctx, query, userID, totalPrice, time.Now()).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("error while creating order: %w", err)
	}

	return orderID, nil
}

func (s *Storage) insertOrderItems(ctx context.Context, tx pgx.Tx, orderID int, items []models.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`

	for _, item := range items {
		_, err := tx.Exec(ctx, query, orderID, item.ProductID, item.Quantity)

		if err != nil {
			return fmt.Errorf("error while inserting order item for product %d: %w", item.ProductID, err)
		}
	}

	return nil
}

func (s *Storage) insertTransaction(ctx context.Context, tx pgx.Tx, orderID int, amount float64) error {
	query := `INSERT INTO transactions (order_id, amount, status, created_at) VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(ctx, query, orderID, amount, "completed", time.Now())
	if err != nil {
		return fmt.Errorf("error while inserting transaction: %w", err)
	}

	return nil
}
