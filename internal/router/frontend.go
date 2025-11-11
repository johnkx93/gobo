package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/app/address"
	"github.com/user/coc/internal/app/order"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/app/user_auth"
	"github.com/user/coc/internal/middleware"
)

// NewFrontendRouter creates the frontend/customer-facing API router
// Routes are prefixed with /api/v1
func NewFrontendRouter(
	userFrontendHandler *user.FrontendHandler,
	orderFrontendHandler *order.FrontendHandler,
	addressFrontendHandler *address.FrontendHandler,
	authHandler *user_auth.Handler,
	authMiddleware func(http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	// Frontend-specific middleware can be added here
	r.Use(middleware.ContentType)

	// Frontend auth routes (public)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
	})

	// Frontend user routes (protected - users can access their own data)
	r.Route("/users", func(r chi.Router) {
		r.Use(authMiddleware)
		// Frontend users can only access/update their own profile
		r.Get("/me", userFrontendHandler.GetMe)    // Get current user profile
		r.Put("/me", userFrontendHandler.UpdateMe) // Update current user profile
	})

	// Frontend order routes (protected - users can manage their own orders)
	r.Route("/orders", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", orderFrontendHandler.CreateOrder)    // Create own order
		r.Get("/", orderFrontendHandler.ListOrders)      // List own orders
		r.Get("/{id}", orderFrontendHandler.GetOrder)    // Get own order by ID
		r.Put("/{id}", orderFrontendHandler.UpdateOrder) // Update own order
		// Note: Delete is not exposed on frontend (business decision)
	})

	// Frontend address routes (protected - users can manage their own addresses)
	r.Route("/addresses", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", addressFrontendHandler.CreateAddress)            // Create own address
		r.Get("/", addressFrontendHandler.ListAddresses)             // List own addresses
		r.Get("/{id}", addressFrontendHandler.GetAddress)            // Get own address by ID
		r.Put("/{id}", addressFrontendHandler.UpdateAddress)         // Update own address
		r.Delete("/{id}", addressFrontendHandler.DeleteAddress)      // Delete own address
		r.Post("/default", addressFrontendHandler.SetDefaultAddress) // Set default address
	})

	return r
}
