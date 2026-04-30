package models

import "github.com/google/uuid"

type Cart struct {
	ID     uuid.UUID  `db:"id" json:"id"`
	UserID uuid.UUID  `db:"user_id" json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CartID    uuid.UUID `db:"cart_id" json:"cart_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`

	Product *Product `json:"product,omitempty"`
}
