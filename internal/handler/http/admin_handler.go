package http

import (
	"encoding/json"
	"net/http"

	models "ecommers/internal/domain"
	admin "ecommers/internal/service/admin"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AdminHandler struct {
	adminService admin.AdminService
}

func NewAdminHandler(service admin.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: service,
	}
}

// RegisterRoutes Registers admin routes with role verification.
func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.With(RequireAdmin).Post("/product", h.CreateProduct)
		r.With(RequireAdmin).Delete("/product/{id}", h.DeleteProduct)
		r.With(RequireAdmin).Get("/orders", h.GetOrders)
		r.With(RequireAdmin).Delete("/user/{id}", h.DeleteUser)
	})
}

// CreateProduct godoc
// @Summary Create a new product
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body models.Product true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/admin/product [post]
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.Product

	// Decoding JSON from the request body into a product structure.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Calling the service
	product, err := h.adminService.CreateProduct(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning a successful response in JSON format.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct godoc
// @Summary Delete product by ID
// @Tags Admin
// @Param id path string true "Product UUID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid UUID format"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/admin/product/{id} [delete]
func (h *AdminHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// We extract the ID from the URL string parameters using chi
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Calling the service
	if err := h.adminService.DeleteProduct(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Administrator access required
	// Returns an empty 204 No Content response.
	w.WriteHeader(http.StatusNoContent)
}

// GetOrders godoc
// @Summary Get all orders
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Order
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/admin/orders [get]
func (h *AdminHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.adminService.AllOrders(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Tags Admin
// @Param id path string true "User UUID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid UUID format"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/admin/user/{id} [delete]
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	if err := h.adminService.DeleteUser(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
