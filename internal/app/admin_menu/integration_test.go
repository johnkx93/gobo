package admin_menu

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/coc/internal/db"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, *db.Queries) {
	t.Helper()

	// Use DATABASE_URL from environment (Docker Postgres)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration tests")
	}

	ctx := context.Background()

	// Parse and create pool config
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("failed to parse database URL: %v", err)
	}

	// Configure pool for testing
	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("failed to create connection pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	// Create queries
	queries := db.New(pool)

	// Cleanup on test completion
	t.Cleanup(func() {
		pool.Close()
	})

	return pool, queries
}

func createTestPermission(t *testing.T, ctx context.Context, queries *db.Queries, code, name, category string) db.Permission {
	t.Helper()
	perm, err := queries.CreatePermission(ctx, db.CreatePermissionParams{
		Code:     code,
		Name:     name,
		Category: category,
	})
	if err != nil {
		t.Fatalf("failed to create test permission: %v", err)
	}
	return perm
}

func createTestMenuItem(t *testing.T, ctx context.Context, queries *db.Queries, parentID, code, label, icon, path string, permissionID pgtype.UUID, orderIndex int32) db.MenuItem {
	t.Helper()
	params := db.CreateMenuItemParams{
		Code:       code,
		Label:      label,
		OrderIndex: orderIndex,
	}

	if parentID != "" {
		parentUUID, err := uuid.Parse(parentID)
		if err != nil {
			t.Fatalf("invalid parent ID: %v", err)
		}
		params.ParentID = pgtype.UUID{Bytes: parentUUID, Valid: true}
	}

	if icon != "" {
		params.Icon = pgtype.Text{String: icon, Valid: true}
	}

	if path != "" {
		params.Path = pgtype.Text{String: path, Valid: true}
	}

	if permissionID.Valid {
		params.PermissionID = permissionID
	}

	menuItem, err := queries.CreateMenuItem(ctx, params)
	if err != nil {
		t.Fatalf("failed to create test menu item: %v", err)
	}
	return menuItem
}

func assignPermissionToRole(t *testing.T, ctx context.Context, queries *db.Queries, role string, permissionID pgtype.UUID) {
	t.Helper()
	err := queries.AssignPermissionToRole(ctx, db.AssignPermissionToRoleParams{
		Role:         role,
		PermissionID: permissionID,
	})
	if err != nil {
		t.Fatalf("failed to assign permission to role: %v", err)
	}
}

func TestIntegration_GetMenuForRole_SuperAdmin(t *testing.T) {
	pool, queries := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)

	// Clean up any existing test data
	_, _ = pool.Exec(ctx, "DELETE FROM role_permissions WHERE role IN ('super_admin', 'admin', 'moderator')")
	_, _ = pool.Exec(ctx, "DELETE FROM menu_items WHERE code LIKE '%super%' OR code LIKE '%admin%' OR code LIKE '%mod%'")
	_, _ = pool.Exec(ctx, "DELETE FROM permissions WHERE code LIKE 'dashboard.%' OR code LIKE 'users.%'")

	// Create test permissions
	dashboardPerm := createTestPermission(t, ctx, qtx, "dashboard.view", "View Dashboard", "dashboard")
	usersPerm := createTestPermission(t, ctx, qtx, "users.manage", "Manage Users", "users")

	// Create menu items (just root items for simplicity)
	_ = createTestMenuItem(t, ctx, qtx, "", "dashboard_super", "Dashboard", "dashboard-icon", "/admin/dashboard", pgtype.UUID{Valid: false}, 1)
	_ = createTestMenuItem(t, ctx, qtx, "", "users_super", "User Management", "users-icon", "/admin/users", usersPerm.ID, 2)

	// Assign permissions to super_admin role
	assignPermissionToRole(t, ctx, qtx, "super_admin", dashboardPerm.ID)
	assignPermissionToRole(t, ctx, qtx, "super_admin", usersPerm.ID)

	// Test GetMenuForRole for super_admin
	menuItems, err := GetMenuForRole(ctx, qtx, "super_admin")
	if err != nil {
		t.Fatalf("GetMenuForRole failed: %v", err)
	}

	// Find our test menu items (they should be present among all menu items)
	var foundDashboard, foundUsers *MenuItem
	for _, item := range menuItems {
		if item.Label == "Dashboard" {
			foundDashboard = item
		}
		if item.Label == "User Management" {
			foundUsers = item
		}
	}

	if foundDashboard == nil {
		t.Error("expected to find Dashboard menu item")
	} else {
		if foundDashboard.Order != 1 {
			t.Errorf("expected dashboard item with order 1, got order %d", foundDashboard.Order)
		}
	}

	if foundUsers == nil {
		t.Error("expected to find User Management menu item")
	} else {
		if foundUsers.Order != 2 {
			t.Errorf("expected users item with order 2, got order %d", foundUsers.Order)
		}
		if foundUsers.Path != "/admin/users" {
			t.Errorf("expected users item with path '/admin/users', got '%s'", foundUsers.Path)
		}
	}
}

func TestIntegration_GetMenuForRole_AdminLimitedPermissions(t *testing.T) {
	pool, queries := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)

	// Clean up any existing test data
	_, _ = pool.Exec(ctx, "DELETE FROM role_permissions WHERE role IN ('super_admin', 'admin', 'moderator')")
	_, _ = pool.Exec(ctx, "DELETE FROM menu_items WHERE code LIKE '%super%' OR code LIKE '%admin%' OR code LIKE '%mod%'")
	_, _ = pool.Exec(ctx, "DELETE FROM permissions WHERE code LIKE 'dashboard.%' OR code LIKE 'users.%'")

	// Create test permissions
	dashboardPerm := createTestPermission(t, ctx, qtx, "dashboard.view", "View Dashboard", "dashboard")
	_ = createTestPermission(t, ctx, qtx, "users.manage", "Manage Users", "users")

	// Create menu items - one with permission, one without
	_ = createTestMenuItem(t, ctx, qtx, "", "dashboard_admin", "Dashboard", "dashboard-icon", "/admin/dashboard", dashboardPerm.ID, 1)
	_ = createTestMenuItem(t, ctx, qtx, "", "settings_admin", "Settings", "settings-icon", "/admin/settings", pgtype.UUID{Valid: false}, 2) // No permission required

	// Only assign dashboard permission to admin role
	assignPermissionToRole(t, ctx, qtx, "admin", dashboardPerm.ID)
	// Don't assign users permission

	// Test GetMenuForRole for admin (should see dashboard and settings among all items)
	menuItems, err := GetMenuForRole(ctx, qtx, "admin")
	if err != nil {
		t.Fatalf("GetMenuForRole failed: %v", err)
	}

	// Should see both dashboard (has permission) and settings (no permission required) among the menu items
	labels := make(map[string]bool)
	for _, item := range menuItems {
		labels[item.Label] = true
	}

	if !labels["Dashboard"] {
		t.Error("expected to find Dashboard item")
	}
	if !labels["Settings"] {
		t.Error("expected to find Settings item")
	}
}

func TestIntegration_GetMenuForRole_ModeratorNoPermissions(t *testing.T) {
	pool, queries := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)

	// Clean up any existing test data
	_, _ = pool.Exec(ctx, "DELETE FROM role_permissions WHERE role IN ('super_admin', 'admin', 'moderator')")
	_, _ = pool.Exec(ctx, "DELETE FROM menu_items WHERE code LIKE '%super%' OR code LIKE '%admin%' OR code LIKE '%mod%'")
	_, _ = pool.Exec(ctx, "DELETE FROM permissions WHERE code LIKE 'dashboard.%' OR code LIKE 'users.%'")

	// Create test permissions
	dashboardPerm := createTestPermission(t, ctx, qtx, "dashboard.view", "View Dashboard", "dashboard")
	settingsPerm := createTestPermission(t, ctx, qtx, "settings.manage", "Manage Settings", "settings")

	// Create menu items that require permissions
	_ = createTestMenuItem(t, ctx, qtx, "", "dashboard_mod", "Dashboard", "dashboard-icon", "/admin/dashboard", dashboardPerm.ID, 1)
	_ = createTestMenuItem(t, ctx, qtx, "", "settings_mod", "Settings", "settings-icon", "/admin/settings", settingsPerm.ID, 2) // Requires permission

	// Don't assign any permissions to moderator role

	// Test GetMenuForRole for moderator (should not see permission-required items)
	menuItems, err := GetMenuForRole(ctx, qtx, "moderator")
	if err != nil {
		t.Fatalf("GetMenuForRole failed: %v", err)
	}

	// Should not see dashboard or settings items since they require permissions that moderator doesn't have
	paths := make(map[string]bool)
	for _, item := range menuItems {
		paths[item.Path] = true
	}

	if paths["/admin/dashboard"] {
		t.Error("should not find /admin/dashboard item (requires permission)")
	}
	if paths["/admin/settings"] {
		t.Error("should not find /admin/settings item (requires permission)")
	}
}

func TestIntegration_GetRolePermissions(t *testing.T) {
	pool, queries := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)

	// Create test permissions
	perm1 := createTestPermission(t, ctx, qtx, "test.perm1", "Test Permission 1", "test")
	perm2 := createTestPermission(t, ctx, qtx, "test.perm2", "Test Permission 2", "test")
	_ = createTestPermission(t, ctx, qtx, "other.perm", "Other Permission", "other")

	// Assign permissions to test role
	assignPermissionToRole(t, ctx, qtx, "test_role", perm1.ID)
	assignPermissionToRole(t, ctx, qtx, "test_role", perm2.ID)
	// Don't assign perm3

	// Test GetRolePermissions
	permissions, err := GetRolePermissions(ctx, qtx, "test_role")
	if err != nil {
		t.Fatalf("GetRolePermissions failed: %v", err)
	}

	if len(permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(permissions))
	}

	if !permissions["test.perm1"] {
		t.Error("expected to find test.perm1 permission")
	}
	if !permissions["test.perm2"] {
		t.Error("expected to find test.perm2 permission")
	}
	if permissions["other.perm"] {
		t.Error("should not find other.perm permission")
	}
}

func TestIntegration_GetRolePermissions_NoPermissions(t *testing.T) {
	pool, queries := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)

	// Create test permissions but don't assign any
	createTestPermission(t, ctx, qtx, "test.perm1", "Test Permission 1", "test")

	// Test GetRolePermissions for role with no permissions
	permissions, err := GetRolePermissions(ctx, qtx, "empty_role")
	if err != nil {
		t.Fatalf("GetRolePermissions failed: %v", err)
	}

	if len(permissions) != 0 {
		t.Errorf("expected 0 permissions for role with no assignments, got %d", len(permissions))
	}
}
