package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/validation"
)

// Helper function to create admin request with admin role in context
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

// Helper function to create request without admin role (for testing auth failures)
func newRequestWithoutAdminRole(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// No admin role in context
	rec := httptest.NewRecorder()
	return req, rec
}

// TestHandler_CreateAdmin_MissingAdminRole tests CreateAdmin without admin role
func TestHandler_CreateAdmin_MissingAdminRole(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req := httptest.NewRequest("POST", "/admins", nil)
	// No admin role in context
	rec := httptest.NewRecorder()

	handler.CreateAdmin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestHandler_GetAdmin_MissingAdminRole tests GetAdmin without admin role
func TestHandler_GetAdmin_MissingAdminRole(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newRequestWithoutAdminRole("GET", "/admins/"+uuid.New().String(), nil)

	handler.GetAdmin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestHandler_ListAdmins_MissingAdminRole tests ListAdmins without admin role
func TestHandler_ListAdmins_MissingAdminRole(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newRequestWithoutAdminRole("GET", "/admins", nil)

	handler.ListAdmins(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestHandler_UpdateAdmin_MissingAdminRole tests UpdateAdmin without admin role
func TestHandler_UpdateAdmin_MissingAdminRole(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newRequestWithoutAdminRole("PUT", "/admins/"+uuid.New().String(), map[string]interface{}{
		"first_name": "Updated",
	})

	handler.UpdateAdmin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestHandler_DeleteAdmin_MissingAdminRole tests DeleteAdmin without admin role
func TestHandler_DeleteAdmin_MissingAdminRole(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newRequestWithoutAdminRole("DELETE", "/admins/"+uuid.New().String(), nil)

	handler.DeleteAdmin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

// TestHandler_GetAdmin_MissingEntityID tests GetAdmin with missing entity ID
func TestHandler_GetAdmin_MissingEntityID(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("GET", "/admins/", nil)

	handler.GetAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_UpdateAdmin_MissingEntityID tests UpdateAdmin with missing entity ID
func TestHandler_UpdateAdmin_MissingEntityID(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/admins/", map[string]interface{}{
		"first_name": "Updated",
	})

	handler.UpdateAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_DeleteAdmin_MissingEntityID tests DeleteAdmin with missing entity ID
func TestHandler_DeleteAdmin_MissingEntityID(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("DELETE", "/admins/", nil)

	handler.DeleteAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_CreateAdmin_InvalidJSON tests CreateAdmin with invalid JSON
func TestHandler_CreateAdmin_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("POST", "/admins", nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

	handler.CreateAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_UpdateAdmin_InvalidJSON tests UpdateAdmin with invalid JSON
func TestHandler_UpdateAdmin_InvalidJSON(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/admins/"+uuid.New().String(), nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

	handler.UpdateAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_CreateAdmin_ValidationError tests CreateAdmin with validation error
func TestHandler_CreateAdmin_ValidationError(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("POST", "/admins", map[string]interface{}{
		"email":    "invalid-email",
		"username": "tu",  // too short
		"password": "123", // too short
	})

	handler.CreateAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// TestHandler_UpdateAdmin_ValidationError tests UpdateAdmin with validation error
func TestHandler_UpdateAdmin_ValidationError(t *testing.T) {
	handler := NewHandler(nil, validation.New())
	req, rec := newAdminRequest("PUT", "/admins/"+uuid.New().String(), map[string]interface{}{
		"email":    "invalid-email",
		"username": "tu", // too short
	})

	handler.UpdateAdmin(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
