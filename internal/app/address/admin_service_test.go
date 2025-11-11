package address

import (
	"context"
	"testing"

	"github.com/google/uuid"
	apperrors "github.com/user/coc/internal/errors"
)

// Test helper functions
func TestGetStringValue(t *testing.T) {
	t.Run("nil pointer returns empty string", func(t *testing.T) {
		result := getStringValue(nil)
		if result != "" {
			t.Errorf("expected empty string, got '%s'", result)
		}
	})

	t.Run("non-nil pointer returns value", func(t *testing.T) {
		value := "test value"
		result := getStringValue(&value)
		if result != "test value" {
			t.Errorf("expected 'test value', got '%s'", result)
		}
	})
}

// Test UUID validation
func TestAdminService_CreateAddress_InvalidUUID(t *testing.T) {
	service := &AdminService{queries: nil, auditService: nil}
	ctx := context.Background()

	req := CreateAddressRequest{
		UserID:  "invalid-uuid",
		Address: "123 Test Street",
		Floor:   "5",
		UnitNo:  "A",
	}

	_, err := service.CreateAddress(ctx, req)

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

func TestAdminService_GetAddress_InvalidUUID(t *testing.T) {
	service := &AdminService{queries: nil, auditService: nil}
	ctx := context.Background()

	_, err := service.GetAddress(ctx, "invalid-uuid")

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

func TestAdminService_UpdateAddress_InvalidUUID(t *testing.T) {
	service := &AdminService{queries: nil, auditService: nil}
	ctx := context.Background()

	newAddr := "456 New St"
	req := UpdateAddressRequest{Address: &newAddr}

	_, err := service.UpdateAddress(ctx, "invalid-uuid", req)

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

func TestAdminService_DeleteAddress_InvalidUUID(t *testing.T) {
	service := &AdminService{queries: nil, auditService: nil}
	ctx := context.Background()

	err := service.DeleteAddress(ctx, "bad-uuid")

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

func TestAdminService_SetDefaultAddress_InvalidUserID(t *testing.T) {
	service := &AdminService{queries: nil, auditService: nil}
	ctx := context.Background()

	err := service.SetDefaultAddress(ctx, "invalid-uuid", uuid.New().String())

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
