package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/validation"
)

// newAdminRequest creates a new HTTP request with admin role context
func newAdminRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add admin role to context
	ctx := context.WithValue(req.Context(), ctxkeys.AdminRoleContextKey, "super_admin")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	return req, rec
}

func TestAdminHandler_CreateUser_MissingAdminRole(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req := httptest.NewRequest("POST", "/users", nil)
	// No admin role in context
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAdminHandler_CreateUser_InvalidJSON(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("POST", "/users", nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

	handler.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_CreateUser_ValidationError(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("POST", "/users", map[string]interface{}{
		"email": "invalid-email",
	})

	handler.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_GetUser_MissingUserID(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("GET", "/users/", nil)
	// Missing user ID in URL params

	handler.GetUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateUser_MissingUserID(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/users/", nil)
	// Missing user ID in URL params

	handler.UpdateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateUser_InvalidJSON(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/users/123", nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

	handler.UpdateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateUser_ValidationError(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/users/123", map[string]interface{}{
		"email": "invalid-email",
	})

	handler.UpdateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestAdminHandler_DeleteUser_MissingUserID(t *testing.T) {
	handler := NewAdminHandler(nil, validation.New())
	req, rec := newAdminRequest("DELETE", "/users/", nil)
	// Missing user ID in URL params

	handler.DeleteUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
