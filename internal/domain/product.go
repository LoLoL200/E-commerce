package models

import (
	"github.com/google/uuid"
)

// Product DB
type Product struct {
	ID          uuid.UUID `db:"id" json:"id"`
	CategoryID  uuid.UUID `db:"category_id" json:"category_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	Quantity    int       `db:"quantity" json:"quantity"`
}
