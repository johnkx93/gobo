package frontend_auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/coc/internal/validation"
)

func TestHandler_Login_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandler_Login_ValidationError(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	reqBody := `{
		"email": "invalid-email",
		"password": "123"
	}`

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandler_Login_ValidRequestStructure(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	reqBody := `{
		"email": "test@example.com",
		"password": "password123"
	}`

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// This will panic at service call since service is nil, but we just want to test JSON parsing and validation
	// In a real scenario, the service would be properly initialized
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil service - this means JSON parsing and validation passed
			return
		}
	}()

	handler.Login(rec, req)

	// If we get here without panic, something is wrong
	t.Error("expected panic due to nil service")
}

func TestHandler_Register_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandler_Register_ValidationError(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	reqBody := `{
		"email": "invalid-email",
		"username": "ab",
		"password": "123"
	}`

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestHandler_Register_ValidRequestStructure(t *testing.T) {
	handler := NewHandler(nil, validation.New())

	reqBody := `{
		"email": "test@example.com",
		"username": "testuser",
		"password": "password123",
		"first_name": "John",
		"last_name": "Doe"
	}`

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// This will panic at service call since service is nil, but we just want to test JSON parsing and validation
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil service - this means JSON parsing and validation passed
			return
		}
	}()

	handler.Register(rec, req)

	// If we get here without panic, something is wrong
	t.Error("expected panic due to nil service")
}
