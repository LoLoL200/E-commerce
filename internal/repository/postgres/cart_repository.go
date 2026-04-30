// package repository

// import (
// 	"context"
// 	"database/sql"
// 	database "ecommers/internal/db"
// 	models "ecommers/internal/domin"
// 	"ecommers/pkg/utils"
// 	"fmt"

// 	"github.com/google/uuid"
// )

// type CartRepository interface {
// 	// cart
// 	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
// 	GetItemByCartAndProduct(ctx context.Context, cartID, productID uuid.UUID) (*models.CartItem, error)
// 	CreateCart(ctx context.Context, cart *models.Cart) error

// 	// items
// 	AddItem(ctx context.Context, item *models.CartItem) error

// 	UpdateItem(ctx context.Context, itemID uuid.UUID, quantity int) error
// 	UpdateItemByUser(ctx context.Context, userID, itemID uuid.UUID, quantity int) error

// 	DeleteItem(ctx context.Context, itemID uuid.UUID) error
// 	DeleteItemByUser(ctx context.Context, userID, itemID uuid.UUID) error

// 	ClearCart(ctx context.Context, cartID uuid.UUID) error
// }
// type cartRepo struct {
// 	db *database.DB
// }

// // GetCartByUserID implements [CartRepository].
// func (c *cartRepo) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
// 	query := `
// 		SELECT id, user_id
// 		FROM carts
// 		WHERE user_id = $1
// 	`

// 	cart := &models.Cart{}

// 	err := c.db.QueryRowContext(ctx, query, userID).Scan(
// 		&cart.ID,
// 		&cart.UserID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, utils.ErrCartNotFound
// 		}
// 		return nil, fmt.Errorf("get cart error: %w", err)
// 	}

// 	// items
// 	itemsQuery := `
// 		SELECT id, cart_id, product_id, quantity
// 		FROM cart_items
// 		WHERE cart_id = $1
// 	`

// 	rows, err := c.db.QueryContext(ctx, itemsQuery, cart.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("get items error: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item models.CartItem

// 		if err := rows.Scan(
// 			&item.ID,
// 			&item.CartID,
// 			&item.ProductID,
// 			&item.Quantity,
// 		); err != nil {
// 			return nil, fmt.Errorf("scan item error: %w", err)
// 		}

// 		cart.Items = append(cart.Items, item)
// 	}

// 	return cart, nil
// }

// // GetItemByCartAndProduct implements [CartRepository].
// func (c *cartRepo) GetItemByCartAndProduct(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (*models.CartItem, error) {
// 	query := `
// 		SELECT id, cart_id, product_id, quantity
// 		FROM cart_items
// 		WHERE cart_id = $1 AND product_id = $2
// 	`

// 	item := &models.CartItem{}

// 	err := c.db.QueryRowContext(ctx, query, cartID, productID).Scan(
// 		&item.ID,
// 		&item.CartID,
// 		&item.ProductID,
// 		&item.Quantity,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, utils.ErrItemNotFound
// 		}
// 		return nil, fmt.Errorf("get item error: %w", err)
// 	}

// 	return item, nil
// }

// // AddItem implements [CartRepository].
// func (c *cartRepo) AddItem(ctx context.Context, item *models.CartItem) error {
// 	query := `
// 		INSERT INTO cart_items (id, cart_id, product_id, quantity)
// 		VALUES ($1, $2, $3, $4)
// 	`

// 	_, err := c.db.ExecContext(ctx, query,
// 		item.ID,
// 		item.CartID,
// 		item.ProductID,
// 		item.Quantity,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("add item error: %w", err)
// 	}

// 	return nil
// }

// // CreateCart implements [CartRepository].
// func (c *cartRepo) CreateCart(ctx context.Context, cart *models.Cart) error {
// 	query := `
// 		INSERT INTO carts (id, user_id)
// 		VALUES ($1, $2)
// 	`

// 	_, err := c.db.ExecContext(ctx, query,
// 		cart.ID,
// 		cart.UserID,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("create cart error: %w", err)
// 	}

// 	return nil
// }

// // UpdateItem implements [CartRepository].
// func (c *cartRepo) UpdateItem(ctx context.Context, itemID uuid.UUID, quantity int) error {
// 	query := `
// 		UPDATE cart_items
// 		SET quantity = $1
// 		WHERE id = $2
// 	`

// 	res, err := c.db.ExecContext(ctx, query, quantity, itemID)
// 	if err != nil {
// 		return fmt.Errorf("update item error: %w", err)
// 	}

// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("item not found")
// 	}

// 	return nil
// }

// // UpdateItemByUser implements [CartRepository].
// func (c *cartRepo) UpdateItemByUser(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, quantity int) error {
// 	query := `
// 		UPDATE cart_items ci
// 		SET quantity = $1
// 		FROM carts c
// 		WHERE ci.id = $2 AND ci.cart_id = c.id AND c.user_id = $3
// 	`

// 	res, err := c.db.ExecContext(ctx, query, quantity, itemID, userID)
// 	if err != nil {
// 		return fmt.Errorf("update by user error: %w", err)
// 	}

// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("item not found or not owned by user")
// 	}

// 	return nil
// }

// // ClearCart implements [CartRepository].
// func (c *cartRepo) ClearCart(ctx context.Context, cartID uuid.UUID) error {
// 	query := `DELETE FROM cart_items WHERE cart_id = $1`

// 	_, err := c.db.ExecContext(ctx, query, cartID)
// 	if err != nil {
// 		return fmt.Errorf("clear cart error: %w", err)
// 	}

// 	return nil
// }

// // DeleteItem implements [CartRepository].
// func (c *cartRepo) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
// 	query := `DELETE FROM cart_items WHERE id = $1`

// 	res, err := c.db.ExecContext(ctx, query, itemID)
// 	if err != nil {
// 		return fmt.Errorf("delete item error: %w", err)
// 	}

// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("item not found")
// 	}

// 	return nil
// }

// // DeleteItemByUser implements [CartRepository].
// func (c *cartRepo) DeleteItemByUser(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
// 	query := `
// 		DELETE FROM cart_items ci
// 		USING carts c
// 		WHERE ci.id = $1
// 		AND ci.cart_id = c.id
// 		AND c.user_id = $2
// 	`

// 	res, err := c.db.ExecContext(ctx, query, itemID, userID)
// 	if err != nil {
// 		return fmt.Errorf("delete by user error: %w", err)
// 	}

// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("item not found or not owned by user")
// 	}

// 	return nil
// }

//	func NewCartRepository(db *database.DB) CartRepository {
//		return &cartRepo{db: db}
//	}
package repository

import (
	"context"
	"database/sql"
	database "ecommers/internal/db"
	models "ecommers/internal/domin" // Проверь опечатку в названии папки (domin -> domain?)
	"ecommers/pkg/utils"
	"fmt"

	"github.com/google/uuid"
)

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddItem(ctx context.Context, userID uuid.UUID, item *models.CartItem) error
	UpdateItemByUser(ctx context.Context, userID, productID uuid.UUID, quantity int) error
	DeleteItemByUser(ctx context.Context, userID, productID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
	GetItemByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*models.CartItem, error)
}

type cartRepo struct {
	db *database.DB
}

func NewCartRepository(db *database.DB) CartRepository {
	return &cartRepo{db: db}
}

// GetCartByUserID теперь собирает корзину "на лету" из таблицы cart_items
func (c *cartRepo) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	query := `
        SELECT id, product_id, quantity
        FROM cart_items
        WHERE user_id = $1
    `

	rows, err := c.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart items error: %w", err)
	}
	defer rows.Close()

	cart := &models.Cart{
		UserID: userID,
		Items:  []models.CartItem{},
	}

	for rows.Next() {
		var item models.CartItem
		item.CartID = userID // В этой схеме CartID по сути равен UserID
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan item error: %w", err)
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

// AddItem теперь использует user_id напрямую
func (c *cartRepo) AddItem(ctx context.Context, userID uuid.UUID, item *models.CartItem) error {
	query := `
        INSERT INTO cart_items (id, user_id, product_id, quantity)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, product_id) 
        DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
    `

	_, err := c.db.ExecContext(ctx, query,
		uuid.New(), // Генерируем новый ID для записи
		userID,
		item.ProductID,
		item.Quantity,
	)

	if err != nil {
		return fmt.Errorf("add item error: %w", err)
	}

	return nil
}

// UpdateItemByUser обновляет количество товара конкретного юзера
func (c *cartRepo) UpdateItemByUser(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	query := `
        UPDATE cart_items
        SET quantity = $1
        WHERE user_id = $2 AND product_id = $3
    `

	res, err := c.db.ExecContext(ctx, query, quantity, userID, productID)
	if err != nil {
		return fmt.Errorf("update item error: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return utils.ErrItemNotFound
	}

	return nil
}

// DeleteItemByUser удаляет конкретный товар из корзины юзера
func (c *cartRepo) DeleteItemByUser(ctx context.Context, userID, productID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`

	res, err := c.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return fmt.Errorf("delete item error: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return utils.ErrItemNotFound
	}

	return nil
}

// ClearCart очищает всю корзину юзера
func (c *cartRepo) ClearCart(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`

	_, err := c.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("clear cart error: %w", err)
	}

	return nil
}

// GetItemByUserAndProduct нужен для проверки, есть ли уже такой товар в корзине
func (c *cartRepo) GetItemByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (*models.CartItem, error) {
	query := `
        SELECT id, product_id, quantity
        FROM cart_items
        WHERE user_id = $1 AND product_id = $2
    `

	item := &models.CartItem{}
	err := c.db.QueryRowContext(ctx, query, userID, productID).Scan(
		&item.ID,
		&item.ProductID,
		&item.Quantity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}
