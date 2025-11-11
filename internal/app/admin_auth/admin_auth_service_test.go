package admin_auth

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
)

// TestAuthService_GenerateToken tests JWT token generation
func TestAuthService_GenerateToken(t *testing.T) {
	service := &AuthService{
		jwtSecret: "test-secret-key",
	}

	// Create test admin
	adminID := uuid.New()
	var pgUUID pgtype.UUID
	pgUUID.Bytes = adminID
	pgUUID.Valid = true

	admin := &db.Admin{
		ID:       pgUUID,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "admin",
		IsActive: true,
	}

	token, err := service.GenerateToken(admin)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Verify token can be parsed
	parsedToken, err := jwt.ParseWithClaims(token, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.jwtSecret), nil
	})

	if err != nil {
		t.Fatalf("failed to parse generated token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*AdminClaims)
	if !ok || !parsedToken.Valid {
		t.Error("generated token is invalid")
	}

	if claims.AdminID != adminID.String() {
		t.Errorf("expected admin ID %s, got %s", adminID.String(), claims.AdminID)
	}

	if claims.Email != admin.Email {
		t.Errorf("expected email %s, got %s", admin.Email, claims.Email)
	}

	if claims.Username != admin.Username {
		t.Errorf("expected username %s, got %s", admin.Username, claims.Username)
	}

	if claims.Role != admin.Role {
		t.Errorf("expected role %s, got %s", admin.Role, claims.Role)
	}

	if claims.Subject != "admin" {
		t.Errorf("expected subject 'admin', got %s", claims.Subject)
	}

	// Check expiration (should be ~24 hours from now)
	expectedExp := time.Now().Add(24 * time.Hour)
	if claims.ExpiresAt.Time.After(expectedExp.Add(time.Minute)) ||
		claims.ExpiresAt.Time.Before(expectedExp.Add(-time.Minute)) {
		t.Errorf("unexpected expiration time: %v", claims.ExpiresAt.Time)
	}
}

// TestAuthService_ValidateToken tests JWT token validation
func TestAuthService_ValidateToken(t *testing.T) {
	service := &AuthService{
		jwtSecret: "test-secret-key",
	}

	// Create test admin
	adminID := uuid.New()
	var pgUUID pgtype.UUID
	pgUUID.Bytes = adminID
	pgUUID.Valid = true

	admin := &db.Admin{
		ID:       pgUUID,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "admin",
		IsActive: true,
	}

	// Generate token
	token, err := service.GenerateToken(admin)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.AdminID != adminID.String() {
		t.Errorf("expected admin ID %s, got %s", adminID.String(), claims.AdminID)
	}

	if claims.Subject != "admin" {
		t.Errorf("expected subject 'admin', got %s", claims.Subject)
	}
}

// TestAuthService_ValidateToken_InvalidToken tests validation of invalid tokens
func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	service := &AuthService{
		jwtSecret: "test-secret-key",
	}

	tests := []struct {
		name  string
		token string
	}{
		{"empty_token", ""},
		{"invalid_token", "invalid.jwt.token"},
		{"wrong_secret", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ValidateToken(tt.token)
			if err == nil {
				t.Error("expected error for invalid token")
			}
		})
	}
}

// TestAuthService_ValidateToken_WrongSubject tests validation of non-admin tokens
func TestAuthService_ValidateToken_WrongSubject(t *testing.T) {
	service := &AuthService{
		jwtSecret: "test-secret-key",
	}

	// Create a token with wrong subject
	claims := &AdminClaims{
		AdminID:  uuid.New().String(),
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user", // Wrong subject
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(service.jwtSecret))

	_, err := service.ValidateToken(tokenString)
	if err == nil {
		t.Error("expected error for non-admin token")
	}

	if err.Error() != "not an admin token" {
		t.Errorf("expected 'not an admin token' error, got: %v", err)
	}
}

// TestAuthService_GetAdminByID_InvalidUUID tests GetAdminByID with invalid UUID
func TestAuthService_GetAdminByID_InvalidUUID(t *testing.T) {
	service := &AuthService{
		queries: nil,
	}

	ctx := context.Background()
	invalidID := "not-a-uuid"

	_, err := service.GetAdminByID(ctx, invalidID)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	// Should get a validation error about invalid UUID format
	if !strings.Contains(err.Error(), "invalid admin ID format") {
		t.Errorf("expected error about invalid admin ID format, got: %v", err)
	}
}

// TestAuthService_CreateAdmin_InvalidInputs tests CreateAdmin validation
func TestAuthService_CreateAdmin_InvalidInputs(t *testing.T) {
	// CreateAdmin calls database methods first, tested in integration tests
	t.Skip("CreateAdmin calls database methods first, tested in integration tests")
}

// TestAuthService_Login_InvalidInputs tests Login validation
func TestAuthService_Login_InvalidInputs(t *testing.T) {
	// Login calls database methods first, tested in integration tests
	t.Skip("Login calls database methods first, tested in integration tests")
}
