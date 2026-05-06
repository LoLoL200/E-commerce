package http

import (
	"ecommers/internal/service/product"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Product Service
type ProductHandler struct {
	productService product.ProductService
}

// New Product
func NewProductHandler(service product.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: service,
	}
}

func (h *ProductHandler) ProductsRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{id}", h.DetailsProduct)
		r.Get("/search", h.SearchProduct)
		r.Get("/categories", h.ListCategory)
	})
}

// List godoc
// @Summary List all products
// @Description Get a paginated list of products with optional filters and sorting
// @Tags products
// @Produce json
// @Param search query string false "Search by name or description"
// @Param order_by query string false "Sort order (e.g. price_asc, price_desc)"
// @Param limit query int false "Limit results (default 30)"
// @Param offset query int false "Offset for pagination"
// @Param category_id query string false "Filter by category UUID"
// @Success 200 {array} models.Product
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/products [get]
// List Product
func (h *ProductHandler) List(write http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	filter := product.ProductFilter{
		Search:  query.Get("search"),
		OrderBy: query.Get("order_by"),
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(write, http.StatusBadRequest, "invalid limit")
			return
		}
		filter.Limit = limit
	} else {
		filter.Limit = 30
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(write, http.StatusBadRequest, "invalid offset")
			return
		}
		filter.Offset = offset
	}

	if categoryID := query.Get("category_id"); categoryID != "" {
		parsedCategoryID, err := uuid.Parse(categoryID)
		if err != nil {
			respondError(write, http.StatusBadRequest, "invalid category_id")
			return
		}
		filter.CategoryID = &parsedCategoryID
	}

	resp, err := h.productService.ListProduct(request.Context(), filter)
	if err != nil {
		fmt.Println("ERROR:", err)
		respondError(write, http.StatusInternalServerError, "failed to list products")
		return
	}

	respondJSON(write, http.StatusOK, resp)
}

// DetailsProduct godoc
// @Summary Get product details
// @Description Get full information about a single product by UUID
// @Tags products
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} models.Product
// @Failure 400 {object} ErrorResponse "Invalid UUID format"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /api/v1/products/{id} [get]
// Details product(search for Id )
func (h *ProductHandler) DetailsProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "product not found")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// SearchProduct godoc
// @Summary Search products
// @Description Search for products with a required search query string
// @Tags products
// @Produce json
// @Param search query string true "Search query string"
// @Param limit query int false "Limit results (default 20)"
// @Param offset query int false "Offset results"
// @Param order_by query string false "Sort order"
// @Success 200 {array} models.Product
// @Failure 400 {object} ErrorResponse "Search query missing or invalid params"
// @Failure 500 {object} ErrorResponse "Search execution failed"
// @Router /api/v1/products/search [get]
// Search all product
func (h *ProductHandler) SearchProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := product.ProductFilter{
		Search:  query.Get("search"),
		OrderBy: query.Get("order_by"),
		Limit:   20,
	}

	if filter.Search == "" {
		respondError(w, http.StatusBadRequest, "search query is required")
		return
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		filter.Offset = offset
	}

	resp, err := h.productService.ListProduct(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "search failed")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ListCategory godoc
// @Summary List products by category
// @Description Get a list of products that belong to a specific category UUID
// @Tags products
// @Produce json
// @Param category query string true "Category UUID"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} models.Product
// @Failure 400 {object} ErrorResponse "Category ID missing or invalid"
// @Failure 500 {object} ErrorResponse "Query failed"
// @Router /api/v1/products/categories [get]
// List product for category
func (h *ProductHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	categoryID := query.Get("category")
	if categoryID == "" {
		respondError(w, http.StatusBadRequest, "category is required")
		return
	}

	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	filter := product.ProductFilter{
		CategoryID: &parsedCategoryID,
		Limit:      100,
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			respondError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		filter.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			respondError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		filter.Offset = offset
	}

	resp, err := h.productService.ListProductCategory(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list by category")
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
