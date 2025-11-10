package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/user/coc/internal/db"
)

// AdminOnly middleware ensures only admin users can access protected routes
// This is a simple implementation that checks if user email is in ADMIN_EMAILS env var
// For production, consider adding a 'role' column to the users table
func AdminOnly(queries *db.Queries) func(http.Handler) http.Handler {
	// Get admin emails from environment variable
	// Format: ADMIN_EMAILS=admin1@example.com,admin2@example.com
	adminEmailsStr := os.Getenv("ADMIN_EMAILS")
	adminEmails := make(map[string]bool)

	if adminEmailsStr != "" {
		emails := strings.Split(adminEmailsStr, ",")
		for _, email := range emails {
			adminEmails[strings.TrimSpace(strings.ToLower(email))] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (set by auth middleware)
			user, ok := r.Context().Value(UserContextKey).(*db.User)
			if !ok || user == nil {
				respondUnauthorized(w, "user not authenticated")
				return
			}

			// Check if user email is in admin list
			userEmail := strings.ToLower(strings.TrimSpace(user.Email))
			if !adminEmails[userEmail] {
				respondForbidden(w, "access denied: admin privileges required")
				return
			}

			// User is admin, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// AdminOnlyWithRole is a future-proof middleware that checks user role from database
// Use this when you add a 'role' column to users table
// For now, it's a placeholder that uses the same logic as AdminOnly
func AdminOnlyWithRole(queries *db.Queries) func(http.Handler) http.Handler {
	return AdminOnly(queries)

	// Future implementation when role column exists:
	/*
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get user from context
				user, ok := r.Context().Value(UserContextKey).(*db.User)
				if !ok || user == nil {
					respondUnauthorized(w, "user not authenticated")
					return
				}

				// Check if user has admin role
				if user.Role != "admin" {
					respondForbidden(w, "access denied: admin privileges required")
					return
				}

				next.ServeHTTP(w, r)
			})
		}
	*/
}

// WithAdminContext adds an "is_admin" flag to the context
// This can be used by handlers to adjust behavior without blocking access
func WithAdminContext(queries *db.Queries) func(http.Handler) http.Handler {
	adminEmailsStr := os.Getenv("ADMIN_EMAILS")
	adminEmails := make(map[string]bool)

	if adminEmailsStr != "" {
		emails := strings.Split(adminEmailsStr, ",")
		for _, email := range emails {
			adminEmails[strings.TrimSpace(strings.ToLower(email))] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAdmin := false

			// Get user from context
			if user, ok := r.Context().Value(UserContextKey).(*db.User); ok && user != nil {
				userEmail := strings.ToLower(strings.TrimSpace(user.Email))
				isAdmin = adminEmails[userEmail]
			}

			// Add is_admin flag to context
			ctx := context.WithValue(r.Context(), contextKey("is_admin"), isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondForbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"status":false,"message":"` + message + `"}`))
}
