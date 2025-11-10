package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/app/admin"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/app/order"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/middleware"
)

// NewAdminRouter creates the admin panel API router
// Routes are prefixed with /api/admin/v1
func NewAdminRouter(
	userAdminHandler *user.AdminHandler,
	orderAdminHandler *order.AdminHandler,
	adminAuthHandler *admin_auth.AuthHandler,
	adminHandler *admin.Handler,
	adminAuthMiddleware func(http.Handler) http.Handler,
) chi.Router {
	r := chi.NewRouter()

	// Admin-specific middleware can be added here
	r.Use(middleware.ContentType)

	// Admin auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", adminAuthHandler.Login)
		// Admin registration might be restricted or different
		// r.Post("/register", adminAuthHandler.Register)
	})

	// Admin user management (protected)
	r.Route("/users", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all admin user routes
		r.Post("/", userAdminHandler.CreateUser)
		r.Get("/", userAdminHandler.ListUsers)
		r.Get("/{id}", userAdminHandler.GetUser)
		r.Put("/{id}", userAdminHandler.UpdateUser)
		r.Delete("/{id}", userAdminHandler.DeleteUser)
	})

	// Admin management (protected) - only super_admin should access these
	r.Route("/admins", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all admin management routes
		r.Post("/", adminHandler.CreateAdmin)
		r.Get("/", adminHandler.ListAdmins)
		r.Get("/{id}", adminHandler.GetAdmin)
		r.Put("/{id}", adminHandler.UpdateAdmin)
		r.Delete("/{id}", adminHandler.DeleteAdmin)
	})

	// Admin order management (protected)
	r.Route("/orders", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all admin order routes
		r.Post("/", orderAdminHandler.CreateOrder)
		r.Get("/", orderAdminHandler.ListOrders)
		r.Get("/{id}", orderAdminHandler.GetOrder)
		r.Put("/{id}", orderAdminHandler.UpdateOrder)
		r.Delete("/{id}", orderAdminHandler.DeleteOrder)
	})

	return r
}
