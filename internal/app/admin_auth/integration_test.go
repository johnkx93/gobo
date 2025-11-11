package admin_auth

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
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

func TestIntegration_AuthService_Login_Success(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create a test admin
	admin, err := service.CreateAdmin(ctx, "logintest@example.com", "logintest", "password123", "Login", "Test", "admin")
	if err != nil {
		t.Fatalf("failed to create test admin: %v", err)
	}

	adminUUID, err := uuid.FromBytes(admin.ID.Bytes[:])
	if err != nil {
		t.Fatalf("failed to convert admin ID to UUID: %v", err)
	}
	adminIDStr := adminUUID.String()

	// Test login with correct credentials
	token, returnedAdmin, err := service.Login(ctx, "logintest", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	returnedAdminUUID, err := uuid.FromBytes(returnedAdmin.ID.Bytes[:])
	if err != nil {
		t.Fatalf("failed to convert returned admin ID to UUID: %v", err)
	}
	returnedAdminIDStr := returnedAdminUUID.String()

	if returnedAdmin.ID != admin.ID {
		t.Errorf("expected admin ID %s, got %s", adminIDStr, returnedAdminIDStr)
	}

	if returnedAdmin.Email != admin.Email {
		t.Errorf("expected email %s, got %s", admin.Email, returnedAdmin.Email)
	}

	// Verify token is valid
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("token validation failed: %v", err)
	}

	if claims.Subject != "admin" {
		t.Errorf("expected subject 'admin', got %s", claims.Subject)
	}
}

func TestIntegration_AuthService_Login_InvalidCredentials(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create a test admin
	_, err = service.CreateAdmin(ctx, "invalidtest@example.com", "invalidtest", "password123", "Invalid", "Test", "admin")
	if err != nil {
		t.Fatalf("failed to create test admin: %v", err)
	}

	// Test login with wrong password
	_, _, err = service.Login(ctx, "invalidtest", "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password")
	}

	// Test login with non-existent username
	_, _, err = service.Login(ctx, "nonexistent", "password123")
	if err == nil {
		t.Error("expected error for non-existent username")
	}
}

func TestIntegration_AuthService_Login_InactiveAdmin(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create a test admin
	admin, err := service.CreateAdmin(ctx, "inactivetest@example.com", "inactivetest", "password123", "Inactive", "Test", "admin")
	if err != nil {
		t.Fatalf("failed to create test admin: %v", err)
	}

	// Deactivate the admin
	updatedAdmin, err := qtx.UpdateAdmin(ctx, db.UpdateAdminParams{
		ID:           admin.ID,
		Email:        admin.Email,
		Username:     admin.Username,
		PasswordHash: admin.PasswordHash,
		FirstName:    admin.FirstName,
		LastName:     admin.LastName,
		Role:         admin.Role,
		IsActive:     false, // Set inactive
	})
	if err != nil {
		t.Fatalf("failed to deactivate admin: %v", err)
	}

	// Verify admin was deactivated
	if updatedAdmin.IsActive {
		t.Error("expected admin to be inactive")
	}

	// Test login with inactive admin
	_, _, err = service.Login(ctx, "inactivetest", "password123")
	if err == nil {
		t.Error("expected error for inactive admin")
	}

	if err.Error() != "invalid username or password" {
		t.Errorf("expected 'invalid username or password' error for inactive admin, got: %v", err)
	}
}

func TestIntegration_AuthService_CreateAdmin_DuplicateEmail(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create first admin
	_, err = service.CreateAdmin(ctx, "duplicate@example.com", "admin1", "password123", "First", "Admin", "admin")
	if err != nil {
		t.Fatalf("failed to create first admin: %v", err)
	}

	// Try to create second admin with same email
	_, err = service.CreateAdmin(ctx, "duplicate@example.com", "admin2", "password123", "Second", "Admin", "admin")
	if err == nil {
		t.Error("expected error for duplicate email")
	}

	if err.Error() != "admin with this email already exists" {
		t.Errorf("expected duplicate email error, got: %v", err)
	}
}

func TestIntegration_AuthService_CreateAdmin_DuplicateUsername(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create first admin
	_, err = service.CreateAdmin(ctx, "admin1@example.com", "duplicateuser", "password123", "First", "Admin", "admin")
	if err != nil {
		t.Fatalf("failed to create first admin: %v", err)
	}

	// Try to create second admin with same username
	_, err = service.CreateAdmin(ctx, "admin2@example.com", "duplicateuser", "password123", "Second", "Admin", "admin")
	if err == nil {
		t.Error("expected error for duplicate username")
	}

	if err.Error() != "admin with this username already exists" {
		t.Errorf("expected duplicate username error, got: %v", err)
	}
}

func TestIntegration_AuthService_GetAdminByID(t *testing.T) {
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
	service := NewAuthService(qtx, "test-jwt-secret")

	// Create a test admin
	admin, err := service.CreateAdmin(ctx, "gettest@example.com", "gettest", "password123", "Get", "Test", "moderator")
	if err != nil {
		t.Fatalf("failed to create test admin: %v", err)
	}

	// Get admin by ID
	adminUUID, err := uuid.FromBytes(admin.ID.Bytes[:])
	if err != nil {
		t.Fatalf("failed to convert admin ID to UUID: %v", err)
	}
	adminIDStr := adminUUID.String()
	result, err := service.GetAdminByID(ctx, adminIDStr)
	if err != nil {
		t.Fatalf("GetAdminByID failed: %v", err)
	}

	if result.ID != adminIDStr {
		t.Errorf("expected ID %s, got %s", adminIDStr, result.ID)
	}

	if result.Email != admin.Email {
		t.Errorf("expected email %s, got %s", admin.Email, result.Email)
	}

	if result.Username != admin.Username {
		t.Errorf("expected username %s, got %s", admin.Username, result.Username)
	}

	if result.Role != admin.Role {
		t.Errorf("expected role %s, got %s", admin.Role, result.Role)
	}
}
