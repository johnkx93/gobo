package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/app/address"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/app/admin_management"
	"github.com/user/coc/internal/app/admin_menu"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/middleware"
)

// NewAdminRouter creates the admin panel API router
// Routes are prefixed with /api/admin/v1
func NewAdminRouter(
	userAdminHandler *user.AdminHandler,
	addressAdminHandler *address.AdminHandler,
	adminAuthHandler *admin_auth.AuthHandler,
	adminHandler *admin_management.Handler,
	menuHandler *admin_menu.Handler,
	adminAuthMiddleware func(http.Handler) http.Handler,
	permissionMiddleware *middleware.PermissionMiddleware,
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

	// Protected routes that require authentication
	r.Group(func(r chi.Router) {
		r.Use(adminAuthMiddleware)

		// Get admin menu based on role
		r.Get("/menu", menuHandler.GetMenu)
	})

	// Admin user management (protected)
	r.Route("/users", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all admin user routes

		// Create requires users.create permission
		r.With(permissionMiddleware.RequirePermission("users.create")).Post("/", userAdminHandler.CreateUser)

		// Read requires users.read permission
		r.With(permissionMiddleware.RequirePermission("users.read")).Get("/", userAdminHandler.ListUsers)
		r.With(permissionMiddleware.RequirePermission("users.read")).Get("/{id}", userAdminHandler.GetUser)

		// Update requires users.update permission
		r.With(permissionMiddleware.RequirePermission("users.update")).Put("/{id}", userAdminHandler.UpdateUser)

		// Delete requires users.delete permission
		r.With(permissionMiddleware.RequirePermission("users.delete")).Delete("/{id}", userAdminHandler.DeleteUser)
	})

	// Admin management (protected) - only super_admin should access these
	r.Route("/admins", func(r chi.Router) {
		r.Use(adminAuthMiddleware)                                     // Protect all admin management routes
		r.Use(permissionMiddleware.RequirePermission("admins.manage")) // Require admins.manage permission

		r.Post("/", adminHandler.CreateAdmin)
		r.Get("/", adminHandler.ListAdmins)
		r.Get("/{id}", adminHandler.GetAdmin)
		r.Put("/{id}", adminHandler.UpdateAdmin)
		r.Delete("/{id}", adminHandler.DeleteAdmin)
	})

	// (orders feature removed)

	// Admin address management (protected)
	r.Route("/addresses", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all admin address routes

		// Create requires addresses.create permission
		r.With(permissionMiddleware.RequirePermission("addresses.create")).Post("/", addressAdminHandler.CreateAddress)

		// Read requires addresses.read permission
		r.With(permissionMiddleware.RequirePermission("addresses.read")).Get("/", addressAdminHandler.ListAllAddresses)
		r.With(permissionMiddleware.RequirePermission("addresses.read")).Get("/{id}", addressAdminHandler.GetAddress)

		// Update requires addresses.update permission
		r.With(permissionMiddleware.RequirePermission("addresses.update")).Put("/{id}", addressAdminHandler.UpdateAddress)

		// Delete requires addresses.delete permission
		r.With(permissionMiddleware.RequirePermission("addresses.delete")).Delete("/{id}", addressAdminHandler.DeleteAddress)
	})

	// Admin user address management (protected)
	r.Route("/users/{user_id}/addresses", func(r chi.Router) {
		r.Use(adminAuthMiddleware) // Protect all user address routes

		// List user's addresses and set default require addresses.read permission
		r.With(permissionMiddleware.RequirePermission("addresses.read")).Get("/", addressAdminHandler.ListAddressesByUser)
		r.With(permissionMiddleware.RequirePermission("addresses.update")).Post("/default", addressAdminHandler.SetDefaultAddress)
	})

	return r
}
