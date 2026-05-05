package product

import (
	"context"
	models "ecommers/internal/domain"
	repository "ecommers/internal/repository/postgres"
	"ecommers/pkg/utils"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type ProductService interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, request UpdateProductRequest, isAdmin bool) (*models.Product, error)
	ListProduct(ctx context.Context, fillter ProductFilter) (*ProductListResponce, error)
	ListProductCategory(ctx context.Context, fillter ProductFilter) (*ProductListResponce, error)
}
type serviceProduct struct {
	productRepo repository.ProductRepository
}

// GetProduct implements [ProductService].
func (s *serviceProduct) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetByIDProduct(ctx, id)
	if err != nil {
		return nil, utils.ErrUserNotfound
	}
	return product, nil
}

// ListProduct implements [ProductService].
func (s *serviceProduct) ListProduct(ctx context.Context, filter ProductFilter) (*ProductListResponce, error) {
	// Дефолтные значения пагинации
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	products, err := s.productRepo.ListProduct(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("list product error: %w", err)
	}

	// Фильтрация по Search
	if filter.Search != "" {
		filtered := make([]*models.Product, 0)
		search := strings.ToLower(filter.Search)
		for _, p := range products {
			if strings.Contains(strings.ToLower(p.Name), search) ||
				strings.Contains(strings.ToLower(p.Description), search) {
				filtered = append(filtered, p)
			}
		}
		products = filtered
	}

	// Фильтрация по Category
	if filter.CategoryID != nil {
		filtered := make([]*models.Product, 0)
		for _, p := range products {
			if p.CategoryID == *filter.CategoryID {
				filtered = append(filtered, p)
			}
		}
		products = filtered
	}

	// Сортировка
	switch filter.OrderBy {
	case "price_asc":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Price < products[j].Price
		})
	case "price_desc":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Price > products[j].Price
		})
	case "name":
		sort.Slice(products, func(i, j int) bool {
			return products[i].Name < products[j].Name
		})
	}

	return &ProductListResponce{
		Products: products,
		Total:    len(products),
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil

}

// ListProductCategory implements [ProductService].
func (s *serviceProduct) ListProductCategory(ctx context.Context, filter ProductFilter) (*ProductListResponce, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.CategoryID == nil || *filter.CategoryID == uuid.Nil {
		return nil, fmt.Errorf("category is required")
	}

	products, err := s.productRepo.ListProductByCategory(ctx, *filter.CategoryID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("list by category error: %w", err)
	}

	if filter.Search != "" {
		filtered := make([]*models.Product, 0)
		search := strings.ToLower(filter.Search)
		for _, p := range products {
			if strings.Contains(strings.ToLower(p.Name), search) ||
				strings.Contains(strings.ToLower(p.Description), search) {
				filtered = append(filtered, p)
			}
		}
		products = filtered
	}

	return &ProductListResponce{
		Products: products,
		Total:    len(products),
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// UpdateProduct implements [ProductService].
func (s *serviceProduct) UpdateProduct(ctx context.Context, id uuid.UUID, request UpdateProductRequest, isAdmin bool) (*models.Product, error) {
	if !isAdmin {
		return nil, fmt.Errorf("permission denied: admin only")
	}

	//// Get an existing product
	product, err := s.productRepo.GetByIDProduct(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update only the passed fields
	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.Category != nil {
		product.CategoryID = *request.Category
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.Price != nil {
		if *request.Price < 0 {
			return nil, fmt.Errorf("price cannot be negative")
		}
		product.Price = *request.Price
	}
	if request.Quantity != nil {
		product.Quantity = *request.Quantity
	}

	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("update product error: %w", err)
	}

	return product, nil
}

func NewService(productRepo repository.ProductRepository) ProductService {
	return &serviceProduct{productRepo: productRepo}
}
