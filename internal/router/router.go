package router

import (
	"net/http"
	"time"

	"github.com/user/coc/internal/app/order"
	"github.com/user/coc/internal/app/user"
	"github.com/user/coc/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// New creates a new HTTP router with all routes configured
func New(userHandler *user.Handler, orderHandler *order.Handler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentType)

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})

		// Order routes
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", orderHandler.CreateOrder)
			r.Get("/", orderHandler.ListOrders)
			r.Get("/{id}", orderHandler.GetOrder)
			r.Put("/{id}", orderHandler.UpdateOrder)
			r.Delete("/{id}", orderHandler.DeleteOrder)
		})
	})

	return r
}
