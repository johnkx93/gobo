package ctxkeys

import "net/http"

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const (
	// User context keys
	UserContextKey   contextKey = "user"
	UserIDContextKey contextKey = "user_id"

	// Admin context keys
	AdminContextKey     contextKey = "admin"
	AdminIDContextKey   contextKey = "admin_id"
	AdminRoleContextKey contextKey = "admin_role"
)

// GetAdminRole retrieves the admin role from the request context
func GetAdminRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(AdminRoleContextKey).(string)
	return role, ok
}

// GetAdminID retrieves the admin ID from the request context
func GetAdminID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(AdminIDContextKey).(string)
	return id, ok
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(UserIDContextKey).(string)
	return id, ok
}
