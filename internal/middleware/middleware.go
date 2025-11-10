package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/user/coc/internal/app/user_auth"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"
)

// Middleware creates a middleware that validates JWT tokens and adds user to context
func Middleware(authService *user_auth.Service, queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, "missing authorization header")
				return
			}

			// Check if token starts with "Bearer "
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || (parts[0] != "Bearer" && parts[0] != "bearer") {
				respondUnauthorized(w, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				respondUnauthorized(w, "invalid or expired token")
				return
			}

			// Get user from database to ensure user still exists
			user, err := queries.GetUserByEmail(r.Context(), claims.Email)
			if err != nil {
				respondUnauthorized(w, "user not found")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, &user)
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

			// Add user ID to audit context
			userID, err := uuid.Parse(claims.UserID)
			if err == nil {
				ctx = audit.WithUserID(ctx, userID)
			}

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"status":false,"message":"` + message + `"}`))
}
