package order

import (
	"context"
	models "ecommers/internal/domain"
	repository "ecommers/internal/repository/postgres"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateFromCart(ctx context.Context, userID uuid.UUID) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error
}

type service struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
}

func NewService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository) OrderService {
	return &service{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (s *service) CreateFromCart(ctx context.Context, userID uuid.UUID) (*models.Order, error) {

	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cannot create order: cart is empty")
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, ci := range cart.Items {
		itemPrice := 100.0

		orderItems = append(orderItems, models.OrderItem{
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			Price:     itemPrice,
		})
		totalAmount += itemPrice * float64(ci.Quantity)
	}

	newOrder := &models.Order{
		ID:              uuid.New(),
		UserID:          userID,
		Status:          string(models.StatusPending),
		TotalAmount:     totalAmount,
		ShippingAddress: "{}",
		PaymentMethod:   "card",
		Items:           orderItems,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.orderRepo.CreateOrder(ctx, newOrder); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	_ = s.cartRepo.ClearCart(ctx, userID)

	return newOrder, nil
}

func (s *service) GetOrderByID(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("access denied: not your order")
	}

	return order, nil
}

func (s *service) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return s.orderRepo.ListUserOrders(ctx, userID)
}

func (s *service) CancelOrder(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return errors.New("cannot cancel someone else's order")
	}

	return s.orderRepo.UpdateOrderStatus(ctx, orderID, string(models.StatusCancelled))
}
