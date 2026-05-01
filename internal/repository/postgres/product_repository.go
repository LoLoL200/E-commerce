package repository

import (
	"context"
	database "ecommers/internal/db"
	models "ecommers/internal/domin"
	"fmt"

	"github.com/google/uuid"
)

type ProductRepository interface {
	//CreateProduct(ctx context.Context, product *models.Product) error
	GetByIDProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	//GetNameProduct(ctx context.Context, name string) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	//DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProduct(ctx context.Context, limit, offset int) ([]*models.Product, error)
	ListProductByCategory(ctx context.Context, category uuid.UUID, limit, offset int) ([]*models.Product, error)
}

type productRepo struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepo{db: db}
}

// CreateProduct
func (p *productRepo) CreateProduct(ctx context.Context, product *models.Product) error {
	// Вставляем именно в products!
	// Убедись, что имена колонок совпадают с твоей таблицей продуктов
	query := `
        INSERT INTO products (name, description, price, quantity, category_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := p.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Quantity,
		product.CategoryID,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByIDProduct
func (p *productRepo) GetByIDProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
        SELECT id, name, description, price, stock, category_id
        FROM products
        WHERE id = $1`

	product := &models.Product{}
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CategoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// GetNameProduct
func (p *productRepo) GetNameProduct(ctx context.Context, name string) (*models.Product, error) {
	query := `
        SELECT id, name, description, price, quantity, category
        FROM products
        WHERE name = $1`

	product := &models.Product{}
	err := p.db.QueryRowContext(ctx, query, name).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Quantity,
		&product.CategoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}
func (p *productRepo) ListProductByCategory(ctx context.Context, category_id uuid.UUID, limit, offset int) ([]*models.Product, error) {
	query := `
        SELECT id, name, description, price, stock, category_id
        FROM products
      	WHERE category_id = $1
        ORDER BY name
        LIMIT $2 OFFSET $3`

	rows, err := p.db.QueryContext(ctx, query, category_id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list by category error: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity,
			&product.CategoryID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}

// UpdateProduct
func (p *productRepo) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
        UPDATE products
        SET name = $1, description = $2, price = $3, quantity = $4, category = $5
        WHERE id = $6`

	result, err := p.db.ExecContext(ctx, query,
		product.Name,
		product.CategoryID,
		product.Description,
		product.Price,
		product.Quantity,
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

// DeleteProduct
func (p *productRepo) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
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

// ListProduct
func (p *productRepo) ListProduct(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `
        SELECT id, name, description, price, stock, category_id 
        FROM products
        ORDER BY name
        LIMIT $1 OFFSET $2`

	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list products error: %v", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		// Порядок в Scan должен быть ТАКИМ ЖЕ, как в SELECT выше
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Quantity, // Это упадет в поле stock в БД
			&product.CategoryID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		products = append(products, product)
	}
	return products, nil
}
