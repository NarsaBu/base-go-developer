package models

import "time"

type OrderDetail struct {
	OrderID           int
	UserID            int
	TotalPrice        float64
	CreatedAt         time.Time
	Items             []OrderItemDetail
	TransactionID     int
	TransactionStatus string
}

type OrderItemDetail struct {
	ProductID   int
	ProductName string
	Quantity    int
	Price       float64
}

type PopularProduct struct {
	ProductID    int
	ProductName  string
	TotalSold    int
	TotalRevenue float64
}
