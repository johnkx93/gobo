package admin_auth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/coc/internal/validation"
)

// TestAuthHandler_Login_InvalidJSON tests Login with invalid JSON
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	handler := NewAuthHandler(nil, validation.New())
	req := httptest.NewRequest("POST", "/auth/login", nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestAuthHandler_Login_ValidationError tests Login with validation error
func TestAuthHandler_Login_ValidationError(t *testing.T) {
	handler := NewAuthHandler(nil, validation.New())
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{
		"username": "ab",
		"password": "123"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestAuthHandler_Login_ValidRequest tests Login with valid request structure
func TestAuthHandler_Login_ValidRequest(t *testing.T) {
	handler := NewAuthHandler(nil, validation.New())
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{
		"username": "testuser",
		"password": "password123"
	}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// This will panic because service is nil, but we're testing that validation passes
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic due to nil service, but no panic occurred - validation should have passed")
		}
	}()

	handler.Login(rec, req)

	// If we get here, validation passed and service was called
	t.Log("Validation passed, service call attempted")
}
