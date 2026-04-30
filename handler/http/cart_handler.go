package http

import (
	cart "ecommers/internal/service/cart"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartService cart.CartService
}

func NewCartHandler(service cart.CartService) *CartHandler {
	return &CartHandler{
		cartService: service,
	}
}

func (h *CartHandler) RegisterRoutes(r chi.Router) {
	r.Route("/cart", func(r chi.Router) {
		r.Get("/", h.GetCart)

		r.Post("/items", h.AddItem)

		r.Put("/items/{id}", h.UpdateItem)

		r.Delete("/items/{id}", h.RemoveItem)

		r.Delete("/", h.ClearCart)
	})
}

// GetCart godoc
// @Summary Get current user cart
// @Description Retrieve all items in the user's shopping cart
// @Tags cart
// @Produce json
// @Success 200 {object} models.Cart
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cartData, err := h.cartService.GetCart(r.Context(), userID)
	if err != nil {
		// --- ВОТ ЭТО САМОЕ ВАЖНОЕ ---
		log.Printf("🔴 CART SERVICE ERROR: %v", err)
		// ----------------------------
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartData)
}

// GetCart godoc
// @Summary Get current user cart
// @Description Retrieve all items in the user's shopping cart
// @Tags cart
// @Produce json
// @Success 200 {object} models.Cart
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/cart [get]

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG: Entering AddItem handler") // Добавь это

	userID, err := getUserID(r)
	if err != nil {
		log.Printf("DEBUG: getUserID failed: %v", err) // И это
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.cartService.AddItem(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateItem godoc
// @Summary Update item quantity
// @Description Change the quantity of an existing item in the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Item UUID"
// @Param request body UpdateItemRequest true "New quantity"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/cart/items/{id} [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.cartService.UpdateItem(r.Context(), userID, itemID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Delete a specific item from the cart using its ID
// @Tags cart
// @Produce json
// @Param id path string true "Item UUID"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/cart/items/{id} [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	err = h.cartService.RemoveItem(r.Context(), userID, itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ClearCart godoc
// @Summary Clear cart
// @Description Remove all items from the current user's cart
// @Tags cart
// @Produce json
// @Success 200 "OK"
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/cart [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.cartService.ClearCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(ContextKeyUserID)
	if val == nil {
		return uuid.Nil, errors.New("userID not found in context")
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid userID type in context")
	}
	return id, nil
}
