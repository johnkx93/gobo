package admin

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, *db.Queries, *audit.Service) {
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

	// Create queries and services
	queries := db.New(pool)
	auditService := audit.NewService(queries)

	// Cleanup on test completion
	t.Cleanup(func() {
		pool.Close()
	})

	return pool, queries, auditService
}

func TestIntegration_AdminService_CreateAdmin(t *testing.T) {
	pool, queries, auditService := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx) // CRITICAL: Cleanup test data

	// Use transaction for queries
	qtx := queries.WithTx(tx)
	service := NewService(qtx, auditService)

	// Create test data
	req := CreateAdminRequest{
		Email:     "testadmin@example.com",
		Username:  "testadmin",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "Admin",
		Role:      "admin",
	}

	result, err := service.CreateAdmin(ctx, req)
	if err != nil {
		t.Fatalf("CreateAdmin failed: %v", err)
	}

	// Verify data in database
	dbAdmin, err := qtx.GetAdminByID(ctx, stringToPgUUID(result.ID))
	if err != nil {
		t.Fatalf("failed to get admin from database: %v", err)
	}

	if dbAdmin.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, dbAdmin.Email)
	}

	if dbAdmin.Username != req.Username {
		t.Errorf("expected username %s, got %s", req.Username, dbAdmin.Username)
	}

	if dbAdmin.Role != req.Role {
		t.Errorf("expected role %s, got %s", req.Role, dbAdmin.Role)
	}

	if !dbAdmin.IsActive {
		t.Error("expected admin to be active")
	}

	// Verify audit log created
	auditLogs, err := qtx.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: "admins",
		EntityID:   stringToPgUUID(result.ID),
		Limit:      10,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}

	if len(auditLogs) == 0 {
		t.Error("expected audit log entry, got none")
	}

	if auditLogs[0].Action != "CREATE" {
		t.Errorf("expected audit action CREATE, got %s", auditLogs[0].Action)
	}

	// Transaction rolls back automatically - no cleanup needed!
}

func TestIntegration_AdminService_GetAdmin(t *testing.T) {
	pool, queries, auditService := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)
	service := NewService(qtx, auditService)

	// First create an admin
	req := CreateAdminRequest{
		Email:     "gettest@example.com",
		Username:  "gettest",
		Password:  "password123",
		FirstName: "Get",
		LastName:  "Test",
		Role:      "moderator",
	}

	created, err := service.CreateAdmin(ctx, req)
	if err != nil {
		t.Fatalf("failed to create admin for test: %v", err)
	}

	// Now get the admin
	result, err := service.GetAdmin(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetAdmin failed: %v", err)
	}

	if result.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, result.ID)
	}

	if result.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, result.Email)
	}

	if result.Username != req.Username {
		t.Errorf("expected username %s, got %s", req.Username, result.Username)
	}
}

func TestIntegration_AdminService_ListAdmins(t *testing.T) {
	pool, queries, auditService := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)
	service := NewService(qtx, auditService)

	// Create multiple admins
	admins := []CreateAdminRequest{
		{
			Email:     "list1@example.com",
			Username:  "list1",
			Password:  "password123",
			FirstName: "List",
			LastName:  "One",
			Role:      "admin",
		},
		{
			Email:     "list2@example.com",
			Username:  "list2",
			Password:  "password123",
			FirstName: "List",
			LastName:  "Two",
			Role:      "moderator",
		},
	}

	for _, req := range admins {
		_, err := service.CreateAdmin(ctx, req)
		if err != nil {
			t.Fatalf("failed to create admin for list test: %v", err)
		}
	}

	// List admins
	result, err := service.ListAdmins(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListAdmins failed: %v", err)
	}

	if len(result) < 2 {
		t.Errorf("expected at least 2 admins, got %d", len(result))
	}

	// Verify we can find our created admins
	found := 0
	for _, admin := range result {
		if admin.Email == "list1@example.com" || admin.Email == "list2@example.com" {
			found++
		}
	}

	if found != 2 {
		t.Errorf("expected to find 2 created admins, found %d", found)
	}
}

func TestIntegration_AdminService_UpdateAdmin(t *testing.T) {
	pool, queries, auditService := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)
	service := NewService(qtx, auditService)

	// Create an admin
	req := CreateAdminRequest{
		Email:     "updatetest@example.com",
		Username:  "updatetest",
		Password:  "password123",
		FirstName: "Update",
		LastName:  "Test",
		Role:      "moderator",
	}

	created, err := service.CreateAdmin(ctx, req)
	if err != nil {
		t.Fatalf("failed to create admin for update test: %v", err)
	}

	// Update the admin
	updateReq := UpdateAdminRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Role:      "admin",
	}

	result, err := service.UpdateAdmin(ctx, created.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdateAdmin failed: %v", err)
	}

	if result.FirstName != "Updated" {
		t.Errorf("expected first name 'Updated', got '%s'", result.FirstName)
	}

	if result.LastName != "Name" {
		t.Errorf("expected last name 'Name', got '%s'", result.LastName)
	}

	if result.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", result.Role)
	}

	// Verify audit log created
	auditLogs, err := qtx.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: "admins",
		EntityID:   stringToPgUUID(created.ID),
		Limit:      10,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}

	// Should have CREATE and UPDATE logs
	if len(auditLogs) < 2 {
		t.Errorf("expected at least 2 audit logs, got %d", len(auditLogs))
	}

	// Most recent should be UPDATE
	if auditLogs[0].Action != "UPDATE" {
		t.Errorf("expected most recent audit action UPDATE, got %s", auditLogs[0].Action)
	}
}

func TestIntegration_AdminService_DeleteAdmin(t *testing.T) {
	pool, queries, auditService := setupTestDB(t)

	ctx := context.Background()

	// Start transaction for isolation
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Use transaction for queries
	qtx := queries.WithTx(tx)
	service := NewService(qtx, auditService)

	// Create an admin
	req := CreateAdminRequest{
		Email:     "deletetest@example.com",
		Username:  "deletetest",
		Password:  "password123",
		FirstName: "Delete",
		LastName:  "Test",
		Role:      "admin",
	}

	created, err := service.CreateAdmin(ctx, req)
	if err != nil {
		t.Fatalf("failed to create admin for delete test: %v", err)
	}

	// Delete the admin
	err = service.DeleteAdmin(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteAdmin failed: %v", err)
	}

	// Verify admin is marked as inactive (soft delete)
	// GetAdmin should fail since it filters out inactive admins
	_, err = service.GetAdmin(ctx, created.ID)
	if err == nil {
		t.Error("expected GetAdmin to fail for soft-deleted admin")
	}

	// Verify audit log created
	auditLogs, err := qtx.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: "admins",
		EntityID:   stringToPgUUID(created.ID),
		Limit:      10,
		Offset:     0,
	})
	if err != nil {
		t.Fatalf("failed to get audit logs: %v", err)
	}

	// Should have CREATE and DELETE logs
	if len(auditLogs) < 2 {
		t.Errorf("expected at least 2 audit logs, got %d", len(auditLogs))
	}

	// Most recent should be DELETE
	if auditLogs[0].Action != "DELETE" {
		t.Errorf("expected most recent audit action DELETE, got %s", auditLogs[0].Action)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// Helper function to convert string ID to pgtype.UUID
func stringToPgUUID(idStr string) pgtype.UUID {
	id, _ := uuid.Parse(idStr)
	var pgUUID pgtype.UUID
	pgUUID.Bytes = id
	pgUUID.Valid = true
	return pgUUID
}