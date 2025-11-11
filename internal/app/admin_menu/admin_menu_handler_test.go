package admin_menu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/coc/internal/ctxkeys"
)

func newAdminMenuRequest(method, url string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	return req, rec
}

func TestHandler_GetMenu_MissingAdminRole(t *testing.T) {
	// Use nil for queries - this will cause database errors but we want to test context validation first
	handler := NewHandler(nil)
	req, rec := newAdminMenuRequest("GET", "/admin/menu")
	// No admin role in context

	handler.GetMenu(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["message"] != "admin role not found" {
		t.Errorf("expected 'admin role not found' message, got %s", response["message"])
	}
}

func TestHandler_GetMenu_EmptyAdminRole(t *testing.T) {
	// Use nil for queries - this will cause database errors but we want to test context validation first
	handler := NewHandler(nil)
	req, rec := newAdminMenuRequest("GET", "/admin/menu")

	// Add empty admin role to context
	ctx := context.WithValue(req.Context(), ctxkeys.AdminRoleContextKey, "")
	req = req.WithContext(ctx)

	handler.GetMenu(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["message"] != "admin role not found" {
		t.Errorf("expected 'admin role not found' message, got %s", response["message"])
	}
}
