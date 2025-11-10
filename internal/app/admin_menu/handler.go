package admin_menu

import (
	"log/slog"
	"net/http"

	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/response"
)

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
	role, ok := ctxkeys.GetAdminRole(r)
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

	slog.Info("returning menu", "role", role, "root_items", len(menuItems))
	for i, item := range menuItems {
		slog.Info("root menu item", "index", i, "label", item.Label, "children_count", len(item.Children))
	}

	response.JSON(w, http.StatusOK, "menu retrieved successfully", map[string]interface{}{
		"menu": menuItems,
		"role": role,
	})
}
