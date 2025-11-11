package user

import (
	"context"
	"testing"

	"github.com/user/coc/internal/errors"
)

func TestFrontendService_GetUser_InvalidUUID(t *testing.T) {
	service := &FrontendService{
		adminService: nil, // Will be set to nil to test validation
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

func TestFrontendService_UpdateUser_InvalidUUID(t *testing.T) {
	service := &FrontendService{
		adminService: nil, // Will be set to nil to test validation
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
