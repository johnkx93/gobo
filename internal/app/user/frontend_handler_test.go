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

// newUserRequest creates a new HTTP request with user ID context
func newUserRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Add user ID to context
	userID := "550e8400-e29b-41d4-a716-446655440000" // mock UUID
	ctx := context.WithValue(req.Context(), ctxkeys.UserIDContextKey, userID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	return req, rec
}

func TestFrontendHandler_GetMe_MissingUserID(t *testing.T) {
	handler := NewFrontendHandler(nil, validation.New())
	req := httptest.NewRequest("GET", "/users/me", nil)
	// No user ID in context
	rec := httptest.NewRecorder()

	handler.GetMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_UpdateMe_MissingUserID(t *testing.T) {
	handler := NewFrontendHandler(nil, validation.New())
	req := httptest.NewRequest("PUT", "/users/me", nil)
	// No user ID in context
	rec := httptest.NewRecorder()

	handler.UpdateMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestFrontendHandler_UpdateMe_InvalidJSON(t *testing.T) {
	handler := NewFrontendHandler(nil, validation.New())
	req, rec := newUserRequest("PUT", "/users/me", nil)
	req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

	handler.UpdateMe(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestFrontendHandler_UpdateMe_ValidationError(t *testing.T) {
	handler := NewFrontendHandler(nil, validation.New())
	req, rec := newUserRequest("PUT", "/users/me", map[string]interface{}{
		"email": "invalid-email",
	})

	handler.UpdateMe(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}