package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"time"
)

func (s *Storage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	const fn = "storage.postgres.order.CreateOrder"

	query := `INSERT INTO orders (user_id, total_price, created_at) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := s.db.QueryRow(ctx, query, order.UserID, 0, time.Now()).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	const fn = "storage.postgres.order.AddOrderItem"

	query := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order

	err := s.db.QueryRow(ctx, `SELECT id, user_id, total_price, created_at FROM orders WHERE id = $1`, id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt)
	if err != nil {
		return order, fmt.Errorf("%s: %w", fn, err)
	}

	return order, nil
}

func (s *Storage) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	query := `
		SELECT o.id, o.user_id, o.total_price, o.created_at
		FROM orders o
		JOIN users u
			ON u.id = o.user_id
		WHERE u.email = $1
		ORDER BY o.id
`

	rows, err := s.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalPrice, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	query := `
		SELECT id, order_id, product_id, quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
`

	rows, err := s.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var oi models.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orderItems = append(orderItems, oi)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orderItems, nil
}
