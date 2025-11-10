package menu

import (
	"log/slog"
	"net/http"

	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/response"
)

// Context key type for admin context
type contextKey string

const AdminRoleContextKey contextKey = "admin_role"

// Handler handles menu-related requests
type Handler struct {
	queries *db.Queries
}

// NewHandler creates a new menu handler
func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}

// GetMenu returns the menu structure for the authenticated admin
func (h *Handler) GetMenu(w http.ResponseWriter, r *http.Request) {
	// Get admin role from context (set by AdminAuthMiddleware)
	role, ok := r.Context().Value(AdminRoleContextKey).(string)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	// Get menu from database for this role
	menuItems, err := GetMenuForRole(r.Context(), h.queries, role)
	if err != nil {
		slog.Error("failed to get menu", "error", err, "role", role)
		response.Error(w, http.StatusInternalServerError, "failed to retrieve menu")
		return
	}

	response.JSON(w, http.StatusOK, "menu retrieved successfully", map[string]interface{}{
		"menu": menuItems,
		"role": role,
	})
}
