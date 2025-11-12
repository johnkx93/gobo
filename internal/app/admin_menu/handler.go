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
// @Summary      Get admin menu
// @Description  Retrieve the menu structure for the authenticated admin based on their role
// @Tags         Admin Menu
// @Accept       json
// @Produce      json
// @Success      200 {object} response.JSONResponse "Menu retrieved successfully"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Failure      500 {object} response.JSONResponse "Internal server error"
// @Security     BearerAuth
// @Router       /api/admin/v1/menu [get]
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

	response.JSON(w, http.StatusOK, "menu retrieved successfully", map[string]interface{}{
		"menu": menuItems,
		"role": role,
	})
}
