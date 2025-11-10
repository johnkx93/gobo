package middleware

import (
	"log/slog"
	"net/http"

	"github.com/user/coc/internal/app/admin_menu"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/response"
)

// PermissionMiddleware wraps permission checking with database access
type PermissionMiddleware struct {
	queries *db.Queries
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(queries *db.Queries) *PermissionMiddleware {
	return &PermissionMiddleware{
		queries: queries,
	}
}

// RequirePermission middleware ensures the admin has a specific permission
func (pm *PermissionMiddleware) RequirePermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get admin role from context (set by AdminAuthMiddleware)
			role, ok := ctxkeys.GetAdminRole(r)
			if !ok || role == "" {
				response.Error(w, http.StatusUnauthorized, "admin role not found in context")
				return
			}

			// Get permissions for this role from database
			permissions, err := admin_menu.GetRolePermissions(r.Context(), pm.queries, role)
			if err != nil {
				slog.Error("failed to get role permissions", "error", err, "role", role)
				response.Error(w, http.StatusInternalServerError, "failed to check permissions")
				return
			}

			// Check if admin has the required permission
			if !permissions[requiredPermission] {
				slog.Warn("insufficient permissions", "role", role, "required", requiredPermission)
				response.Error(w, http.StatusForbidden, "insufficient permissions for this action")
				return
			}

			// Permission granted, proceed
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission middleware ensures the admin has at least one of the specified permissions
func (pm *PermissionMiddleware) RequireAnyPermission(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := ctxkeys.GetAdminRole(r)
			if !ok || role == "" {
				response.Error(w, http.StatusUnauthorized, "admin role not found in context")
				return
			}

			permissions, err := admin_menu.GetRolePermissions(r.Context(), pm.queries, role)
			if err != nil {
				slog.Error("failed to get role permissions", "error", err, "role", role)
				response.Error(w, http.StatusInternalServerError, "failed to check permissions")
				return
			}

			// Check if admin has any of the required permissions
			hasPermission := false
			for _, perm := range requiredPermissions {
				if permissions[perm] {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				response.Error(w, http.StatusForbidden, "insufficient permissions for this action")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAllPermissions middleware ensures the admin has all of the specified permissions
func (pm *PermissionMiddleware) RequireAllPermissions(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := ctxkeys.GetAdminRole(r)
			if !ok || role == "" {
				response.Error(w, http.StatusUnauthorized, "admin role not found in context")
				return
			}

			permissions, err := admin_menu.GetRolePermissions(r.Context(), pm.queries, role)
			if err != nil {
				slog.Error("failed to get role permissions", "error", err, "role", role)
				response.Error(w, http.StatusInternalServerError, "failed to check permissions")
				return
			}

			// Check if admin has all required permissions
			for _, perm := range requiredPermissions {
				if !permissions[perm] {
					response.Error(w, http.StatusForbidden, "insufficient permissions for this action")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
