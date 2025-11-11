package frontend_auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, *db.Queries) {
	t.Helper()

	// Use DATABASE_URL from environment (Docker Postgres)
	databaseURL := "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable"
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

func TestIntegration_Login_Success(t *testing.T) {
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

	// Create test user
	userID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	testUser := db.User{
		ID:           pgtype.UUID{Bytes: userID, Valid: true},
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		FirstName:    pgtype.Text{String: "John", Valid: true},
		LastName:     pgtype.Text{String: "Doe", Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        testUser.Email,
		Username:     testUser.Username,
		PasswordHash: testUser.PasswordHash,
		FirstName:    testUser.FirstName,
		LastName:     testUser.LastName,
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Test login
	token, user, err := service.Login(ctx, "test@example.com", "testpassword")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", user.Username)
	}

	// Verify token is valid
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("expected email test@example.com in token, got %s", claims.Email)
	}
}

func TestIntegration_Login_InvalidCredentials(t *testing.T) {
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

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Test login with non-existent user
	_, _, err = service.Login(ctx, "nonexistent@example.com", "password")
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	// Test login with wrong password
	// First create a user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_, err = qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	_, _, err = service.Login(ctx, "test@example.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestIntegration_Register_Success(t *testing.T) {
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

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Test registration
	user, err := service.Register(ctx, "newuser@example.com", "newuser", "password123", "Jane", "Smith")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user.Email != "newuser@example.com" {
		t.Errorf("expected email newuser@example.com, got %s", user.Email)
	}

	if user.Username != "newuser" {
		t.Errorf("expected username newuser, got %s", user.Username)
	}

	if user.FirstName.String != "Jane" {
		t.Errorf("expected first name Jane, got %s", user.FirstName.String)
	}

	if user.LastName.String != "Smith" {
		t.Errorf("expected last name Smith, got %s", user.LastName.String)
	}

	// Verify user was created in database
	dbUser, err := qtx.GetUserByEmail(ctx, "newuser@example.com")
	if err != nil {
		t.Fatalf("failed to get user from database: %v", err)
	}

	if dbUser.Username != "newuser" {
		t.Errorf("expected username newuser in database, got %s", dbUser.Username)
	}

	// Verify password was hashed
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte("password123")); err != nil {
		t.Error("password was not properly hashed")
	}
}

func TestIntegration_Register_DuplicateEmail(t *testing.T) {
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

	// Create existing user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_, err = qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        "existing@example.com",
		Username:     "existinguser",
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		t.Fatalf("failed to create existing user: %v", err)
	}

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Try to register with same email
	_, err = service.Register(ctx, "existing@example.com", "newuser", "password123", "John", "Doe")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestIntegration_Register_DuplicateUsername(t *testing.T) {
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

	// Create existing user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_, err = qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        "existing@example.com",
		Username:     "existinguser",
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		t.Fatalf("failed to create existing user: %v", err)
	}

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Try to register with same username
	_, err = service.Register(ctx, "new@example.com", "existinguser", "password123", "John", "Doe")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}

func TestIntegration_TokenGenerationAndValidation(t *testing.T) {
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

	// Create test user
	userID := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	testUser := db.User{
		ID:           pgtype.UUID{Bytes: userID, Valid: true},
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = qtx.CreateUser(ctx, db.CreateUserParams{
		Email:        testUser.Email,
		Username:     testUser.Username,
		PasswordHash: testUser.PasswordHash,
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create service
	auditService := audit.NewService(qtx)
	service := NewService(qtx, auditService, "test-secret", time.Hour)

	// Generate token
	token, err := service.GenerateToken(&testUser)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Verify claims
	if claims.UserID != userID.String() {
		t.Errorf("expected user ID %s, got %s", userID.String(), claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", claims.Email)
	}

	if claims.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", claims.Username)
	}
}
