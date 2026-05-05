package product

import (
	models "ecommers/internal/domain"

	"github.com/google/uuid"
)

// Fillter
type ProductFilter struct {
	Search     string
	CategoryID *uuid.UUID
	Limit      int
	Offset     int
	OrderBy    string
}

// List
type ProductListResponce struct {
	Products []*models.Product
	Total    int
	Limit    int
	Offset   int
}

// Update
type UpdateProductRequest struct {
	Name        *string
	Category    *uuid.UUID
	Description *string
	Price       *float64
	Quantity    *int
}
