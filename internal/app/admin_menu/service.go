package admin_menu

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/coc/internal/db"
)

// MenuItem represents a menu item in the admin panel
type MenuItem struct {
	ID       string      `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Label    string      `json:"label" example:"User Management"`
	Icon     string      `json:"icon,omitempty" example:"users"`
	Path     string      `json:"path,omitempty" example:"/admin/users"`
	Children []*MenuItem `json:"children,omitempty"`
	Order    int         `json:"order" example:"1"`
}

// GetMenuForRole returns the menu structure from database based on admin role
func GetMenuForRole(ctx context.Context, queries *db.Queries, role string) ([]*MenuItem, error) {
	// Fetch menu items from database for this role
	dbMenuItems, err := queries.GetMenuItemsByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	// Build hierarchical menu structure
	return buildMenuTree(dbMenuItems), nil
}

// buildMenuTree converts flat menu items into hierarchical structure
func buildMenuTree(items []db.GetMenuItemsByRoleRow) []*MenuItem {
	// Map to store menu items by ID for quick lookup
	itemMap := make(map[uuid.UUID]*MenuItem)
	var rootItems []*MenuItem

	// First pass: create all menu items
	for _, item := range items {
		itemID, _ := uuid.FromBytes(item.ID.Bytes[:])

		menuItem := &MenuItem{
			ID:       itemID.String(),
			Label:    item.Label,
			Icon:     item.Icon.String,
			Path:     item.Path.String,
			Order:    int(item.OrderIndex),
			Children: []*MenuItem{},
		}

		itemMap[itemID] = menuItem
	}

	// Second pass: build parent-child relationships
	for _, item := range items {
		itemID, _ := uuid.FromBytes(item.ID.Bytes[:])
		menuItem := itemMap[itemID]

		if item.ParentID.Valid {
			// This is a child item
			parentID, _ := uuid.FromBytes(item.ParentID.Bytes[:])
			if parent, exists := itemMap[parentID]; exists {
				parent.Children = append(parent.Children, menuItem)
			}
		} else {
			// This is a root item
			rootItems = append(rootItems, menuItem)
		}
	}

	return rootItems
}

// GetRolePermissions returns permission codes for a given role from database
func GetRolePermissions(ctx context.Context, queries *db.Queries, role string) (map[string]bool, error) {
	permissionCodes, err := queries.GetRolePermissionCodes(ctx, role)
	if err != nil {
		return nil, err
	}

	permissions := make(map[string]bool)
	for _, code := range permissionCodes {
		permissions[code] = true
	}

	return permissions, nil
}
