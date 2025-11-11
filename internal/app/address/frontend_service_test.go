package address

import (
	"context"
	"testing"

	"github.com/google/uuid"
	apperrors "github.com/user/coc/internal/errors"
)

// Test FrontendService UUID validation (only tests validation logic, not DB operations)

func TestFrontendService_GetAddress_InvalidUUID(t *testing.T) {
	service := &FrontendService{queries: nil, auditService: nil}
	ctx := context.Background()
	userID := uuid.New()

	_, err := service.GetAddress(ctx, userID, "invalid-uuid")

	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	domainErr, ok := err.(*apperrors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != apperrors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

func TestFrontendService_UpdateAddress_InvalidUUID(t *testing.T) {
	service := &FrontendService{queries: nil, auditService: nil}
	ctx := context.Background()
	userID := uuid.New()

	newAddr := "456 New St"
	req := UserUpdateAddressRequest{Address: &newAddr}

	_, err := service.UpdateAddress(ctx, userID, "invalid-uuid", req)

	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	domainErr, ok := err.(*apperrors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != apperrors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

func TestFrontendService_DeleteAddress_InvalidUUID(t *testing.T) {
	service := &FrontendService{queries: nil, auditService: nil}
	ctx := context.Background()
	userID := uuid.New()

	err := service.DeleteAddress(ctx, userID, "bad-uuid")

	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	domainErr, ok := err.(*apperrors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != apperrors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}

func TestFrontendService_SetDefaultAddress_InvalidUUID(t *testing.T) {
	service := &FrontendService{queries: nil, auditService: nil}
	ctx := context.Background()
	userID := uuid.New()

	err := service.SetDefaultAddress(ctx, userID, "invalid-uuid")

	if err == nil {
		t.Fatal("expected error for invalid UUID, got nil")
	}

	domainErr, ok := err.(*apperrors.DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}

	if domainErr.Code != apperrors.CodeValidation {
		t.Errorf("expected validation error, got: %s", domainErr.Code)
	}
}
