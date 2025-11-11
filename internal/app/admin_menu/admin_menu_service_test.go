package admin_menu

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
)

func TestBuildMenuTree_EmptyInput(t *testing.T) {
	items := []db.GetMenuItemsByRoleRow{}
	result := buildMenuTree(items)

	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d items", len(result))
	}
}

func TestBuildMenuTree_SingleRootItem(t *testing.T) {
	rootID := uuid.New()
	items := []db.GetMenuItemsByRoleRow{
		{
			ID: pgtype.UUID{Bytes: rootID, Valid: true},
			ParentID: pgtype.UUID{Valid: false}, // No parent
			Code: "dashboard",
			Label: "Dashboard",
			Icon: pgtype.Text{String: "dashboard-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/dashboard", Valid: true},
			PermissionID: pgtype.UUID{Valid: false}, // No permission required
			OrderIndex: 1,
		},
	}

	result := buildMenuTree(items)

	if len(result) != 1 {
		t.Fatalf("expected 1 root item, got %d", len(result))
	}

	item := result[0]
	if item.ID != rootID.String() {
		t.Errorf("expected ID %s, got %s", rootID.String(), item.ID)
	}
	if item.Label != "Dashboard" {
		t.Errorf("expected label 'Dashboard', got %s", item.Label)
	}
	if item.Icon != "dashboard-icon" {
		t.Errorf("expected icon 'dashboard-icon', got %s", item.Icon)
	}
	if item.Path != "/admin/dashboard" {
		t.Errorf("expected path '/admin/dashboard', got %s", item.Path)
	}
	if item.Order != 1 {
		t.Errorf("expected order 1, got %d", item.Order)
	}
	if len(item.Children) != 0 {
		t.Errorf("expected no children, got %d", len(item.Children))
	}
}

func TestBuildMenuTree_ParentChildRelationship(t *testing.T) {
	parentID := uuid.New()
	childID := uuid.New()

	items := []db.GetMenuItemsByRoleRow{
		{
			ID: pgtype.UUID{Bytes: parentID, Valid: true},
			ParentID: pgtype.UUID{Valid: false}, // Root item
			Code: "users",
			Label: "User Management",
			Icon: pgtype.Text{String: "users-icon", Valid: true},
			Path: pgtype.Text{Valid: false}, // No path for parent
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 2,
		},
		{
			ID: pgtype.UUID{Bytes: childID, Valid: true},
			ParentID: pgtype.UUID{Bytes: parentID, Valid: true}, // Child of parent
			Code: "user-list",
			Label: "User List",
			Icon: pgtype.Text{String: "list-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/users", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 1,
		},
	}

	result := buildMenuTree(items)

	if len(result) != 1 {
		t.Fatalf("expected 1 root item, got %d", len(result))
	}

	parent := result[0]
	if parent.ID != parentID.String() {
		t.Errorf("expected parent ID %s, got %s", parentID.String(), parent.ID)
	}
	if parent.Label != "User Management" {
		t.Errorf("expected parent label 'User Management', got %s", parent.Label)
	}
	if len(parent.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(parent.Children))
	}

	child := parent.Children[0]
	if child.ID != childID.String() {
		t.Errorf("expected child ID %s, got %s", childID.String(), child.ID)
	}
	if child.Label != "User List" {
		t.Errorf("expected child label 'User List', got %s", child.Label)
	}
	if child.Path != "/admin/users" {
		t.Errorf("expected child path '/admin/users', got %s", child.Path)
	}
	if len(child.Children) != 0 {
		t.Errorf("expected child to have no children, got %d", len(child.Children))
	}
}

func TestBuildMenuTree_MultipleRootsAndChildren(t *testing.T) {
	// Create multiple root items and children
	dashboardID := uuid.New()
	usersID := uuid.New()
	userListID := uuid.New()
	userCreateID := uuid.New()
	settingsID := uuid.New()

	items := []db.GetMenuItemsByRoleRow{
		// Dashboard (root)
		{
			ID: pgtype.UUID{Bytes: dashboardID, Valid: true},
			ParentID: pgtype.UUID{Valid: false},
			Code: "dashboard",
			Label: "Dashboard",
			Icon: pgtype.Text{String: "dashboard-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/dashboard", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 1,
		},
		// Users (root)
		{
			ID: pgtype.UUID{Bytes: usersID, Valid: true},
			ParentID: pgtype.UUID{Valid: false},
			Code: "users",
			Label: "User Management",
			Icon: pgtype.Text{String: "users-icon", Valid: true},
			Path: pgtype.Text{Valid: false},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 2,
		},
		// User List (child of Users)
		{
			ID: pgtype.UUID{Bytes: userListID, Valid: true},
			ParentID: pgtype.UUID{Bytes: usersID, Valid: true},
			Code: "user-list",
			Label: "User List",
			Icon: pgtype.Text{String: "list-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/users", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 1,
		},
		// User Create (child of Users)
		{
			ID: pgtype.UUID{Bytes: userCreateID, Valid: true},
			ParentID: pgtype.UUID{Bytes: usersID, Valid: true},
			Code: "user-create",
			Label: "Create User",
			Icon: pgtype.Text{String: "plus-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/users/create", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 2,
		},
		// Settings (root)
		{
			ID: pgtype.UUID{Bytes: settingsID, Valid: true},
			ParentID: pgtype.UUID{Valid: false},
			Code: "settings",
			Label: "Settings",
			Icon: pgtype.Text{String: "settings-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/settings", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 3,
		},
	}

	result := buildMenuTree(items)

	if len(result) != 3 {
		t.Fatalf("expected 3 root items, got %d", len(result))
	}

	// Check that items are in correct order
	if result[0].Order != 1 || result[0].Label != "Dashboard" {
		t.Errorf("expected first item to be Dashboard with order 1, got %s with order %d", result[0].Label, result[0].Order)
	}
	if result[1].Order != 2 || result[1].Label != "User Management" {
		t.Errorf("expected second item to be User Management with order 2, got %s with order %d", result[1].Label, result[1].Order)
	}
	if result[2].Order != 3 || result[2].Label != "Settings" {
		t.Errorf("expected third item to be Settings with order 3, got %s with order %d", result[2].Label, result[2].Order)
	}

	// Check Users children
	usersItem := result[1]
	if len(usersItem.Children) != 2 {
		t.Fatalf("expected Users to have 2 children, got %d", len(usersItem.Children))
	}

	// Children should be ordered by order_index
	if usersItem.Children[0].Order != 1 || usersItem.Children[0].Label != "User List" {
		t.Errorf("expected first child to be User List with order 1, got %s with order %d", usersItem.Children[0].Label, usersItem.Children[0].Order)
	}
	if usersItem.Children[1].Order != 2 || usersItem.Children[1].Label != "Create User" {
		t.Errorf("expected second child to be Create User with order 2, got %s with order %d", usersItem.Children[1].Label, usersItem.Children[1].Order)
	}
}

func TestBuildMenuTree_OrphanedChild(t *testing.T) {
	// Test case where a child references a non-existent parent
	childID := uuid.New()
	nonExistentParentID := uuid.New()

	items := []db.GetMenuItemsByRoleRow{
		{
			ID: pgtype.UUID{Bytes: childID, Valid: true},
			ParentID: pgtype.UUID{Bytes: nonExistentParentID, Valid: true}, // Parent doesn't exist
			Code: "orphan",
			Label: "Orphan Child",
			Icon: pgtype.Text{String: "orphan-icon", Valid: true},
			Path: pgtype.Text{String: "/admin/orphan", Valid: true},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 1,
		},
	}

	result := buildMenuTree(items)

	// Orphaned child should not appear in root items
	if len(result) != 0 {
		t.Errorf("expected no root items for orphaned child, got %d", len(result))
	}
}

func TestBuildMenuTree_InvalidUUID(t *testing.T) {
	// Test with invalid UUID bytes (should not panic)
	items := []db.GetMenuItemsByRoleRow{
		{
			ID: pgtype.UUID{Bytes: [16]byte{}, Valid: true}, // Invalid UUID
			ParentID: pgtype.UUID{Valid: false},
			Code: "test",
			Label: "Test Item",
			Icon: pgtype.Text{Valid: false},
			Path: pgtype.Text{Valid: false},
			PermissionID: pgtype.UUID{Valid: false},
			OrderIndex: 1,
		},
	}

	// This should not panic
	result := buildMenuTree(items)

	if len(result) != 1 {
		t.Fatalf("expected 1 item despite invalid UUID, got %d", len(result))
	}

	// The ID will be the zero UUID string representation
	if result[0].ID == "" {
		t.Error("expected non-empty ID string")
	}
}