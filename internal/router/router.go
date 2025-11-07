package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/app/auth"
	"github.com/user/coc/internal/app/order"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/middleware"
)

// New creates a new HTTP router with all routes configured
func New(userHandler *user.Handler, orderHandler *order.Handler, authHandler *auth.Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
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

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentType)

		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
		})

		// User routes (protected)
		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware) // Protect all user routes
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Order routes (protected)
		r.Route("/orders", func(r chi.Router) {
			r.Use(authMiddleware) // Protect all order routes
			r.Post("/", orderHandler.CreateOrder)
			r.Get("/", orderHandler.ListOrders)
			r.Get("/{id}", orderHandler.GetOrder)
			r.Put("/{id}", orderHandler.UpdateOrder)
			r.Delete("/{id}", orderHandler.DeleteOrder)
		})
	})

	return r
}
