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

// Helper to create request with user context
func newUserRequest(method, url string, body interface{}) *http.Request {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context (simulating auth middleware)
	ctx := context.WithValue(req.Context(), ctxkeys.UserIDContextKey, uuid.New().String())
	return req.WithContext(ctx)
}

// Test CreateAddress handler
func TestFrontendHandler_CreateAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	// Request WITHOUT user ID in context
	req := httptest.NewRequest(http.MethodPost, "/api/v1/addresses", nil)
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
}

func TestFrontendHandler_CreateAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := newUserRequest(http.MethodPost, "/api/v1/addresses", nil)
	req.Body = http.NoBody
	rec := httptest.NewRecorder()

	handler.CreateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestFrontendHandler_CreateAddress_ValidationError(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	// Missing required fields
	invalidReq := UserCreateAddressRequest{}

	req := newUserRequest(http.MethodPost, "/api/v1/addresses", invalidReq)
	rec := httptest.NewRecorder()

	handler.CreateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test GetAddress handler
func TestFrontendHandler_GetAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.GetAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_GetAddress_MissingAddressID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := newUserRequest(http.MethodGet, "/api/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.GetAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test UpdateAddress handler
func TestFrontendHandler_UpdateAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.UpdateAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_UpdateAddress_MissingAddressID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := newUserRequest(http.MethodPut, "/api/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.UpdateAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestFrontendHandler_UpdateAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	addressID := uuid.New().String()
	req := newUserRequest(http.MethodPut, "/api/v1/addresses/"+addressID, nil)
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
func TestFrontendHandler_DeleteAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/addresses/123", nil)
	rec := httptest.NewRecorder()

	handler.DeleteAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_DeleteAddress_MissingAddressID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := newUserRequest(http.MethodDelete, "/api/v1/addresses/", nil)
	rec := httptest.NewRecorder()

	handler.DeleteAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

// Test ListAddresses handler
func TestFrontendHandler_ListAddresses_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/addresses", nil)
	rec := httptest.NewRecorder()

	handler.ListAddresses(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

// Test SetDefaultAddress handler
func TestFrontendHandler_SetDefaultAddress_MissingUserID(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/addresses/default", nil)
	rec := httptest.NewRecorder()

	handler.SetDefaultAddress(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_SetDefaultAddress_InvalidJSON(t *testing.T) {
	validator := validation.New()
	handler := NewFrontendHandler(nil, validator)

	req := newUserRequest(http.MethodPost, "/api/v1/addresses/default", nil)
	req.Body = http.NoBody
	rec := httptest.NewRecorder()

	handler.SetDefaultAddress(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}
