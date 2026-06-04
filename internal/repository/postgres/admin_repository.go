package repository

import (
	"context"
	database "ecommers/internal/db"
	models "ecommers/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

// Admin Repository
type AdminRepositorty interface {
	// Order
	ListOrder(ctx context.Context) ([]*models.Order, error)
	RedactStatus(ctx context.Context, status *models.Order) error
	//Product Управління товарами: створення, редагування, видалення, зміна залишків
	CreateProduct(ctx context.Context, product *models.Product) error
	RedactProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	// User
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUser(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type adminRepo struct {
	db *database.DB
}

// CreateProduct implements [AdminRepositorty].
func (a *adminRepo) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
        INSERT INTO products (name, description, price, stock, category_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := a.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// Delete implements [AdminRepositorty].

func (a *adminRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := a.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete product error: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// RedactProduct implements [AdminRepositorty].
func (a *adminRepo) RedactProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET
			name = $1,
			description = $2,
			price = $3,
			update_at = NOW()
		WHERE id = $4	

	`

	_, err := a.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct implements [AdminRepositorty].
func (a *adminRepo) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, stock = $4, category = $5
        WHERE id = $6`

	result, err := a.db.ExecContext(ctx, query,
		product.Name,
		product.CategoryID,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("update product error: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// List implements [AdminRepositorty].
func (a *adminRepo) ListUser(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	query := `
        SELECT id, user_id,	status, total_anmount, shipping_adress, payment_mathod, created_at, updated_at 
        FROM products
        ORDER BY name
        LIMIT $1 OFFSET $2`

	rows, err := a.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list products error: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.Surname,
			&user.Role,
			&user.CreateAt,
			&user.UpdateAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		users = append(users, user)
	}
	return users, nil

}

// Delete user or bloled/unblocked implements [AdminRepositorty].
func (a *adminRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user WHERE user_id = $1`

	_, err := a.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("clear cart error: %w", err)
	}

	return nil
}

// Update implements [AdminRepositorty].
func (a *adminRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET  email= $1, password_hash = $2, first_name = $3, surname = $4, role = $5
        WHERE id = $6`

	result, err := a.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.Surname,
		user.Role,
	)
	if err != nil {
		return fmt.Errorf("update product error: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// ListOrder implements [AdminRepositorty].
func (a *adminRepo) ListOrder(ctx context.Context) ([]*models.Order, error) {
	query := `
       SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at 
        FROM orders
`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list orders error: %v", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan error: %v", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// RedactStatus implements [AdminRepositorty].
func (a *adminRepo) RedactStatus(ctx context.Context, status *models.Order) error {
	panic("unimplemented")
}

func NewAdminRepository(db *database.DB) AdminRepositorty {
	return &adminRepo{db: db}
}
