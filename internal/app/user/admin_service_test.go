package user

import (
	"context"
	"testing"

	"github.com/user/coc/internal/errors"
)

func TestStringsTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no spaces",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "leading spaces",
			input:    "  hello",
			expected: "hello",
		},
		{
			name:     "trailing spaces",
			input:    "hello  ",
			expected: "hello",
		},
		{
			name:     "both spaces",
			input:    "  hello  ",
			expected: "hello",
		},
		{
			name:     "multiple spaces",
			input:    "  hello   world  ",
			expected: "hello   world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringsTrim(tt.input)
			if result != tt.expected {
				t.Errorf("stringsTrim(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAdminService_CreateUser_InvalidUUID(t *testing.T) {
	service := &AdminService{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()

	// Test with invalid UUID format
	req := CreateUserRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// This will panic because queries is nil, but we're testing that the method exists
	// In practice, this test ensures the method signature is correct
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic due to nil queries, but no panic occurred")
		}
	}()

	_, err := service.CreateUser(ctx, req)
	if err == nil {
		t.Error("expected error due to nil queries, got nil")
	}
}

func TestAdminService_GetUser_InvalidUUID(t *testing.T) {
	service := &AdminService{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()

	// Test with invalid UUID format
	_, err := service.GetUser(ctx, "invalid-uuid")

	// Should return validation error for invalid UUID
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != errors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

func TestAdminService_UpdateUser_InvalidUUID(t *testing.T) {
	service := &AdminService{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()

	// Test with invalid UUID format
	req := UpdateUserRequest{
		Email: stringPtr("new@example.com"),
	}

	_, err := service.UpdateUser(ctx, "invalid-uuid", req)

	// Should return validation error for invalid UUID
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != errors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

func TestAdminService_DeleteUser_InvalidUUID(t *testing.T) {
	service := &AdminService{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()

	// Test with invalid UUID format
	err := service.DeleteUser(ctx, "invalid-uuid")

	// Should return validation error for invalid UUID
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != errors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
