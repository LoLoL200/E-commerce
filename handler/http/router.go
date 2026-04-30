package http

import (
	auth "ecommers/internal/service/auth"
	userService "ecommers/internal/service/auth"
	cart "ecommers/internal/service/cart"
	"ecommers/internal/service/order"
	"ecommers/internal/service/product"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterConfig struct {
	AuthService    auth.AuthService
	UserService    userService.UserService
	ProductService product.ProductService
	CartService    cart.CartService
	OrderService   order.OrderService
}

func NewRouter(config RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	r.Route("/api/v1", func(r chi.Router) {

		// ======================
		// AUTH (public)
		// ======================
		authHandler := NewAuthHandler(config.AuthService)
		authHandler.RegisterRoutes(r)

		// ======================
		// PRODUCT (public)
		// ======================
		productHandler := NewProductHandler(config.ProductService)
		productHandler.ProductsRoutes(r)

		// ======================
		// PROTECTED ROUTES (Auth Required)
		// ======================
		r.Group(func(r chi.Router) {
			r.Use(RequireAuth(config.AuthService))

			// CART
			cartHandler := NewCartHandler(config.CartService)
			cartHandler.RegisterRoutes(r)

			// USERS
			userHandler := NewUserHandler(config.UserService, config.AuthService)
			userHandler.RegsterRoutes(r)

			// ORDERS (Добавляем сюда)
			// ======================
			// ORDER (protected)
			// ======================
			orderHandler := NewOrderHandler(config.OrderService)
			orderHandler.OrderRoutes(r) // 3. Регистрируем маршруты заказов
		})
	})

	return r
}
