package postgres

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"time"
)

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	query := `
		SELECT 
			o.id AS order_id,
			o.user_id,
			o.total_price,
			o.created_at,
			t.id AS transaction_id,
			t.status AS transaction_status,
			p.id AS product_id,
			p.name AS product_name,
			oi.quantity,
			p.price
		FROM orders o
		JOIN users u ON o.user_id = u.id
		JOIN order_items oi ON o.id = oi.order_id
		JOIN products p ON oi.product_id = p.id
		JOIN transactions t ON o.id = t.order_id
		WHERE u.email = $1
		ORDER BY o.created_at DESC, o.id, oi.id
	`

	rows, err := s.db.Query(ctx, query, email)
	if err != nil {
		return nil, fmt.Errorf("error while querying order history: %w", err)
	}
	defer rows.Close()

	ordersMap := make(map[int]*models.OrderDetail)

	var orderIDs []int

	for rows.Next() {
		var (
			orderID           int
			userID            int
			totalPrice        float64
			createdAt         time.Time
			transactionID     int
			transactionStatus string
			productID         int
			productName       string
			quantity          int
			price             float64
		)

		err := rows.Scan(
			&orderID,
			&userID,
			&totalPrice,
			&createdAt,
			&transactionID,
			&transactionStatus,
			&productID,
			&productName,
			&quantity,
			&price,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning order history: %w", err)
		}

		if _, exists := ordersMap[orderID]; !exists {
			ordersMap[orderID] = &models.OrderDetail{
				OrderID:           orderID,
				UserID:            userID,
				TotalPrice:        totalPrice,
				CreatedAt:         createdAt,
				TransactionID:     transactionID,
				TransactionStatus: transactionStatus,
				Items:             []models.OrderItemDetail{},
			}
			orderIDs = append(orderIDs, orderID)
		}

		ordersMap[orderID].Items = append(ordersMap[orderID].Items, models.OrderItemDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    quantity,
			Price:       price,
		})

		ordersMap[orderID].Items = append(ordersMap[orderID].Items, models.OrderItemDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    quantity,
			Price:       price,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating order history: %w", err)
	}

	orders := make([]models.OrderDetail, 0, len(ordersMap))
	for _, orderID := range orderIDs {
		orders = append(orders, *ordersMap[orderID])
	}

	return orders, nil
}

func (s *Storage) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	query := `
		SELECT 
			p.id AS product_id,
			p.name AS product_name,
			SUM(oi.quantity) AS total_sold,
			SUM(oi.quantity * p.price) AS total_revenue
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC
	`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while querying popular products: %w", err)
	}
	defer rows.Close()

	var products []models.PopularProduct

	for rows.Next() {
		var product models.PopularProduct
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.TotalSold,
			&product.TotalRevenue,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning popular product: %w", err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating popular products: %w", err)
	}

	return products, nil
}
