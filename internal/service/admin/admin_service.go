package admin

import (
	"context"
	models "ecommers/internal/domain"
	repository "ecommers/internal/repository/postgres"
	"fmt"

	"github.com/google/uuid"
)

type AdminService interface {

	//Product
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	//User
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateProfile(ctx context.Context, user *models.User) (*models.User, error)
	//Orders

	AllOrders(ctx context.Context) (*models.ListOrder, error)

	//All

	//Statistic
}

type serviceAdmin struct {
	adminRepo repository.AdminRepositorty
	userRepo  repository.UserRepository
}

// AllOrders implements [AdminService].
func (s *serviceAdmin) AllOrders(ctx context.Context) (*models.ListOrder, error) {

	orders, err := s.adminRepo.ListOrder(ctx)
	if err != nil {
		return nil, fmt.Errorf("List orders error: %w", err)
	}

	return &models.ListOrder{
		Orders:   orders,
		Quantity: len(orders),
	}, nil
}

// CreateProduct implements [AdminService].
func (s *serviceAdmin) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	// Inspection Name
	if product.Name == "" {
		return nil, fmt.Errorf("Description is void")
	}

	// Inspection Description
	if product.Description == "" {
		return nil, fmt.Errorf("Description is void")
	}

	// Inspection Price
	if product.Price <= 0 {
		return nil, fmt.Errorf("Qrice must be greater than zero")
	}

	// Inspection Quantity
	if product.Stock <= 0 {
		return nil, fmt.Errorf("Quantit must be greater than zero")
	}

	err := s.adminRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("service create product error: %w", err)
	}
	return product, nil
}

// DeleteProduct implements [AdminService].
func (s *serviceAdmin) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	err := s.adminRepo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("service delete product error: %w", err)
	}
	return nil
}

// DeleteUser implements [AdminService].
func (s *serviceAdmin) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.adminRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("service delete user error: %w", err)
	}
	return nil
}

// UpdateProfile implements [AdminService].
func (s *serviceAdmin) UpdateProfile(ctx context.Context, user *models.User) (*models.User, error) {
	err := s.adminRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service update user error: %w", err)
	}
	return user, nil
}

// ListProduct(ctx context.Context, fillter ProductFilter) (*ProductListResponce, error)
func NewService(adminRepo repository.AdminRepositorty) AdminService {
	return &serviceAdmin{adminRepo: adminRepo}
}
