package models

import (
	"time"

	"github.com/google/uuid"
)

// Order статус заказа
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID              uuid.UUID   `db:"id" json:"id"`
	UserID          uuid.UUID   `db:"user_id" json:"user_id"`
	Status          string      `db:"status" json:"status"`
	TotalAmount     float64     `db:"total_amount" json:"total_amount"`
	ShippingAddress interface{} `db:"shipping_address" json:"shipping_address"`
	PaymentMethod   string      `db:"payment_method" json:"payment_method"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
	Items           []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	OrderID   uuid.UUID `db:"order_id" json:"order_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	Price     float64   `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
