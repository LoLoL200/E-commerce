package http

import (
	"ecommers/internal/service/order"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService order.OrderService
}

func NewOrderHandler(service order.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: service,
	}
}

func (h *OrderHandler) OrderRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", h.Create)           // Создать заказ
		r.Get("/", h.ListMyOrders)      // Список моих заказов
		r.Get("/{id}", h.DetailsOrder)  // Детали конкретного заказа
		r.Put("/{id}/cancel", h.Cancel) // Отмена заказа
	})
}

// Create godoc
// @Summary Create order from cart
// @Description Creates a new order using all items currently in the user's cart
// @Tags orders
// @Accept json
// @Produce json
// @Success 201 {object} models.Order
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/orders [post]
// POST /api/v1/orders
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: service address: %p\n", h.orderService)
	userID, err := getUserID(r) // Твоя функция извлечения ID из контекста
	if err != nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	order, err := h.orderService.CreateFromCart(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create order: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, order)
}

// ListMyOrders godoc
// @Summary List user orders
// @Description Get a list of all orders belonging to the authenticated user
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders [get]
// GET /api/v1/orders
func (h *OrderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := h.orderService.GetUserOrders(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch orders")
		return
	}

	respondJSON(w, http.StatusOK, orders)
}

// DetailsOrder godoc
// @Summary Get order details
// @Description Get detailed information about a specific order by its ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} models.Order
// @Failure 400 {object} ErrorResponse "Invalid ID format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Router /api/v1/orders/{id} [get]
// GET /api/v1/orders/{id}
func (h *OrderHandler) DetailsOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "order not found or access denied")
		return
	}

	respondJSON(w, http.StatusOK, order)
}

// Cancel godoc
// @Summary Cancel an order
// @Description Set order status to cancelled
// @Tags orders
// @Produce json
// @Param id path string true "Order ID (UUID)"
// @Success 200 {object} map[string]string "{"status": "cancelled"}"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/orders/{id}/cancel [put]
// PUT /api/v1/orders/{id}/cancel
func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	err = h.orderService.CancelOrder(r.Context(), orderID, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to cancel order: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}
