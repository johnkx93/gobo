package frontend_auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
)

func TestService_GenerateToken(t *testing.T) {
	service := &Service{
		jwtSecret:           "test-secret",
		bearerTokenDuration: time.Hour,
	}

	userID := uuid.New()
	user := &db.User{
		ID:       pgtype.UUID{Bytes: userID, Valid: true},
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	// Validate the token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

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

func TestService_ValidateToken_InvalidToken(t *testing.T) {
	service := &Service{
		jwtSecret: "test-secret",
	}

	_, err := service.ValidateToken("invalid.token.here")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestService_ValidateToken_WrongSecret(t *testing.T) {
	service1 := &Service{
		jwtSecret:           "secret1",
		bearerTokenDuration: time.Hour,
	}

	service2 := &Service{
		jwtSecret: "secret2", // Different secret
	}

	user := &db.User{
		ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := service1.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = service2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for token signed with different secret, got nil")
	}
}

func TestService_ValidateToken_ExpiredToken(t *testing.T) {
	service := &Service{
		jwtSecret:           "test-secret",
		bearerTokenDuration: -time.Hour, // Already expired
	}

	user := &db.User{
		ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = service.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestToUserResponse(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	userResp := ToUserResponse(
		pgtype.UUID{Bytes: userID, Valid: true},
		"test@example.com",
		"testuser",
		pgtype.Text{String: "John", Valid: true},
		pgtype.Text{String: "Doe", Valid: true},
		pgtype.Timestamptz{Time: now, Valid: true},
		pgtype.Timestamptz{Time: now, Valid: true},
	)

	if userResp.ID != userID.String() {
		t.Errorf("expected ID %s, got %s", userID.String(), userResp.ID)
	}

	if userResp.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", userResp.Email)
	}

	if userResp.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", userResp.Username)
	}

	if userResp.FirstName == nil || *userResp.FirstName != "John" {
		t.Errorf("expected first name John, got %v", userResp.FirstName)
	}

	if userResp.LastName == nil || *userResp.LastName != "Doe" {
		t.Errorf("expected last name Doe, got %v", userResp.LastName)
	}
}

func TestToUserResponse_InvalidUUID(t *testing.T) {
	userResp := ToUserResponse(
		pgtype.UUID{Valid: false}, // Invalid UUID
		"test@example.com",
		"testuser",
		pgtype.Text{Valid: false},
		pgtype.Text{Valid: false},
		pgtype.Timestamptz{Valid: false},
		pgtype.Timestamptz{Valid: false},
	)

	if userResp.ID != "" {
		t.Errorf("expected empty ID for invalid UUID, got %s", userResp.ID)
	}

	if userResp.FirstName != nil {
		t.Error("expected nil first name for invalid text")
	}

	if userResp.LastName != nil {
		t.Error("expected nil last name for invalid text")
	}
}
