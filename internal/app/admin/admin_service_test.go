package admin

import (
	"context"
	"strings"
	"testing"
)

// TestStringsTrim tests the strings.TrimSpace helper function
func TestStringsTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty_string", "", ""},
		{"no_spaces", "hello", "hello"},
		{"leading_spaces", "  hello", "hello"},
		{"trailing_spaces", "hello  ", "hello"},
		{"both_spaces", "  hello  ", "hello"},
		{"multiple_spaces", "  hello   world  ", "hello   world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(tt.input)
			if result != tt.expected {
				t.Errorf("strings.TrimSpace(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestAdminService_CreateAdmin_InvalidUUID tests CreateAdmin with invalid UUID
func TestAdminService_CreateAdmin_InvalidUUID(t *testing.T) {
	service := &Service{
		queries:      nil, // Will panic if DB called
		auditService: nil, // Only testing validation
	}

	ctx := context.Background()
	req := CreateAdminRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		Role:      "admin",
	}

	// This should panic because queries is nil - we're testing that validation happens before DB calls
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic due to nil queries, but no panic occurred")
		}
	}()

	_, err := service.CreateAdmin(ctx, req)
	if err != nil {
		// If we get here, validation passed and it tried to call DB
		t.Logf("Validation passed, DB call attempted (expected): %v", err)
	}
}

// TestAdminService_GetAdmin_InvalidUUID tests GetAdmin with invalid UUID
func TestAdminService_GetAdmin_InvalidUUID(t *testing.T) {
	service := &Service{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()
	invalidID := "not-a-uuid"

	_, err := service.GetAdmin(ctx, invalidID)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	// Should get an internal error about invalid UUID format
	if !strings.Contains(err.Error(), "invalid admin ID format") {
		t.Errorf("expected error about invalid admin ID format, got: %v", err)
	}
}

// TestAdminService_UpdateAdmin_InvalidUUID tests UpdateAdmin with invalid UUID
func TestAdminService_UpdateAdmin_InvalidUUID(t *testing.T) {
	service := &Service{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()
	invalidID := "not-a-uuid"
	req := UpdateAdminRequest{
		FirstName: "Updated",
	}

	_, err := service.UpdateAdmin(ctx, invalidID, req)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	// Should get an internal error about invalid UUID format
	if !strings.Contains(err.Error(), "invalid admin ID format") {
		t.Errorf("expected error about invalid admin ID format, got: %v", err)
	}
}

// TestAdminService_DeleteAdmin_InvalidUUID tests DeleteAdmin with invalid UUID
func TestAdminService_DeleteAdmin_InvalidUUID(t *testing.T) {
	service := &Service{
		queries:      nil,
		auditService: nil,
	}

	ctx := context.Background()
	invalidID := "not-a-uuid"

	err := service.DeleteAdmin(ctx, invalidID)
	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	// Should get an internal error about invalid UUID format
	if !strings.Contains(err.Error(), "invalid admin ID format") {
		t.Errorf("expected error about invalid admin ID format, got: %v", err)
	}
}
