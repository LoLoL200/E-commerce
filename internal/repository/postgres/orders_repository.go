package repository

import (
	"context"
	database "ecommers/internal/db"
	models "ecommers/internal/domain" // Проверь опечатку в пути (domin -> domain?)
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error
}

type orderRepo struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) OrderRepository {
	return &orderRepo{db: db}
}

// CreateOrder: Now accounts for ALL NOT NULL columns in your schemas
func (r *orderRepo) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	queryOrder := `
        INSERT INTO orders (id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	payMethod := order.PaymentMethod
	if payMethod == "" {
		payMethod = "card"
	}

	_, err = tx.ExecContext(ctx, queryOrder,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		order.ShippingAddress,
		payMethod,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	queryItem := `
        INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryItem,
			uuid.New(),
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
			time.Now(),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetOrderByID: Field mapping corrected.
func (r *orderRepo) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	queryOrder := `
        SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at 
        FROM orders WHERE id = $1`

	order := &models.Order{}
	err := r.db.QueryRowContext(ctx, queryOrder, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.PaymentMethod,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	queryItems := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, queryItems, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// ListUserOrders: Query fixed (total_amount вместо price)
func (r *orderRepo) ListUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	query := `
        SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at 
        FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		o := &models.Order{}
		err := rows.Scan(
			&o.ID, &o.UserID, &o.Status, &o.TotalAmount,
			&o.ShippingAddress, &o.PaymentMethod, &o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// Upadate Order Status
func (r *orderRepo) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}
