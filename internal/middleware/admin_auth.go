package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/response"
)

const (
	AdminContextKey     contextKey = "admin"
	AdminIDContextKey   contextKey = "admin_id"
	AdminRoleContextKey contextKey = "admin_role"
)

// AdminAuthMiddleware checks for valid admin JWT tokens
func AdminAuthMiddleware(authService *admin_auth.AuthService) func(http.Handler) http.Handler {
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

			// Validate admin token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired admin token")
				return
			}

			// Get admin from database to ensure admin still exists and is active
			adminUser, err := authService.GetAdminByID(r.Context(), claims.AdminID)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "admin not found")
				return
			}

			// Check if admin is still active
			if !adminUser.IsActive {
				response.Error(w, http.StatusUnauthorized, "admin account is disabled")
				return
			}

			// Add admin to context
			ctx := context.WithValue(r.Context(), AdminContextKey, &adminUser)
			ctx = context.WithValue(ctx, AdminIDContextKey, claims.AdminID)
			ctx = context.WithValue(ctx, AdminRoleContextKey, claims.Role)

			// Add admin ID to audit context
			adminID, err := uuid.Parse(claims.AdminID)
			if err == nil {
				ctx = audit.WithUserID(ctx, adminID)
			}

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdminRole middleware ensures the admin has specific role(s)
func RequireAdminRole(allowedRoles ...string) func(http.Handler) http.Handler {
	roleMap := make(map[string]bool)
	for _, role := range allowedRoles {
		roleMap[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get admin role from context
			role, ok := r.Context().Value(AdminRoleContextKey).(string)
			if !ok || role == "" {
				response.Error(w, http.StatusForbidden, "admin role not found in context")
				return
			}

			// Check if admin has required role
			if !roleMap[role] {
				response.Error(w, http.StatusForbidden, "insufficient privileges: required role not found")
				return
			}

			// Admin has required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}
