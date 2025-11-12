package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/user/coc/internal/app/address"
	"github.com/user/coc/internal/app/admin"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/app/admin_menu"
	"github.com/user/coc/internal/app/frontend_auth"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/middleware"

	_ "github.com/user/coc/docs/swagger" // Import generated swagger docs
)

// New creates a new HTTP router with all routes configured
// This router mounts both admin and frontend API routers
func New(
	userAdminHandler *user.AdminHandler,
	userFrontendHandler *user.FrontendHandler,
	addressAdminHandler *address.AdminHandler,
	addressFrontendHandler *address.FrontendHandler,
	userAuthHandler *frontend_auth.Handler,
	adminAuthHandler *admin_auth.AuthHandler,
	adminHandler *admin.Handler,
	menuHandler *admin_menu.Handler,
	userAuthMiddleware func(http.Handler) http.Handler,
	adminAuthMiddleware func(http.Handler) http.Handler,
	permissionMiddleware *middleware.PermissionMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)    // Add request ID to all requests
	r.Use(middleware.AuditContext) // Add IP and user agent to all requests

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Swagger UI endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Mount frontend API router (customer-facing)
	r.Mount("/api/v1", NewFrontendRouter(
		userFrontendHandler,
		addressFrontendHandler,
		userAuthHandler,
		userAuthMiddleware,
	))

	// Mount admin API router (admin panel)
	r.Mount("/api/admin/v1", NewAdminRouter(
		userAdminHandler,
		addressAdminHandler,
		adminAuthHandler,
		adminHandler,
		menuHandler,
		adminAuthMiddleware,
		permissionMiddleware,
	))

	return r
}
