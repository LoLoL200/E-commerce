package cart

import (
	"context"
	models "ecommers/internal/domain"
	repository "ecommers/internal/repository/postgres"
	"errors"

	"github.com/google/uuid"
)

type CartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error)
	AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) error // в этой схеме работаем по productID
	RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type service struct {
	cartRepo repository.CartRepository
}

func NewService(cartRepo repository.CartRepository) CartService {
	return &service{cartRepo: cartRepo}
}

// GetCart
func (s *service) GetCart(ctx context.Context, userID uuid.UUID) (*models.Cart, error) {
	return s.cartRepo.GetCartByUserID(ctx, userID)
}

// AddItem
func (s *service) AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	item := &models.CartItem{
		ProductID: productID,
		Quantity:  quantity,
	}

	return s.cartRepo.AddItem(ctx, userID, item)
}

// UpdateItem
func (s *service) UpdateItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	return s.cartRepo.UpdateItemByUser(ctx, userID, productID, quantity)
}

// RemoveItem
func (s *service) RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	if productID == uuid.Nil {
		return errors.New("invalid product id")
	}

	return s.cartRepo.DeleteItemByUser(ctx, userID, productID)
}

// ClearCart
func (s *service) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.cartRepo.ClearCart(ctx, userID)
}
