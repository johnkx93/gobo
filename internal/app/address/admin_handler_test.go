package address

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/validation"
)

// Helper to create request with admin context
func newAdminRequest(method, url string, body interface{}) *http.Request {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add admin role to context using the correct key
	ctx := context.WithValue(req.Context(), ctxkeys.AdminRoleContextKey, "super_admin")
	return req.WithContext(ctx)
}

// Test CreateAddress handler
func TestAdminHandler_CreateAddress_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	// Request WITHOUT admin role in context
	req := httptest.NewRequest(http.MethodPost, "/api/admin/v1/addresses", nil)
	rec := httptest.NewRecorder()

	handler.CreateAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&response)

	if response["status"].(bool) {
		t.Error("expected status false in response")
	}

	if !contains(response["message"].(string), "admin role") {
		t.Errorf("expected admin role error message, got: %s", response["message"])
	}
}

func TestAdminHandler_CreateAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodPost, "/api/admin/v1/addresses", nil)
	req.Body = http.NoBody
	rec := httptest.NewRecorder()

	handler.CreateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestAdminHandler_CreateAddress_ValidationError(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	// Missing required fields
	invalidReq := CreateAddressRequest{}

	req := newAdminRequest(http.MethodPost, "/api/admin/v1/addresses", invalidReq)
	rec := httptest.NewRecorder()

	handler.CreateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test GetAddress handler
func TestAdminHandler_GetAddress_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.GetAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAdminHandler_GetAddress_MissingID(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodGet, "/api/admin/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.GetAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test UpdateAddress handler
func TestAdminHandler_UpdateAddress_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.UpdateAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateAddress_MissingID(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodPut, "/api/admin/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.UpdateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestAdminHandler_UpdateAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	addressID := uuid.New().String()
	req := newAdminRequest(http.MethodPut, "/api/admin/v1/addresses/"+addressID, nil)
	req.Body = http.NoBody

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", addressID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()

	handler.UpdateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test DeleteAddress handler
func TestAdminHandler_DeleteAddress_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.DeleteAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAdminHandler_DeleteAddress_MissingID(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodDelete, "/api/admin/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.DeleteAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test SetDefaultAddress handler
func TestAdminHandler_SetDefaultAddress_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/v1/users/123/addresses/default", nil)
	rec := httptest.NewRecorder()

	handler.SetDefaultAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAdminHandler_SetDefaultAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodPost, "/api/admin/v1/users//addresses/default", nil)
	rec := httptest.NewRecorder()

	handler.SetDefaultAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestAdminHandler_SetDefaultAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	userID := uuid.New().String()
	req := newAdminRequest(http.MethodPost, "/api/admin/v1/users/"+userID+"/addresses/default", nil)
	req.Body = http.NoBody

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("user_id", userID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()

	handler.SetDefaultAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test ListAllAddresses handler
func TestAdminHandler_ListAllAddresses_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/v1/addresses", nil)
	rec := httptest.NewRecorder()

	handler.ListAllAddresses(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

// Test ListAddressesByUser handler
func TestAdminHandler_ListAddressesByUser_MissingAdminRole(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/v1/users/123/addresses", nil)
	rec := httptest.NewRecorder()

	handler.ListAddressesByUser(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestAdminHandler_ListAddressesByUser_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewAdminHandler(nil, validator)

	req := newAdminRequest(http.MethodGet, "/api/admin/v1/users//addresses", nil)
	rec := httptest.NewRecorder()

	handler.ListAddressesByUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != ""
}
