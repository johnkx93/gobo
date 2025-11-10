package user

import (
	"encoding/json"
	"net/http"

	"github.com/user/coc/internal/middleware"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// FrontendHandler handles frontend/customer-facing user operations
// Users can only access and modify their OWN data
type FrontendHandler struct {
	service  *Service
	validate *validation.Validator
}

func NewFrontendHandler(service *Service, validator *validation.Validator) *FrontendHandler {
	return &FrontendHandler{
		service:  service,
		validate: validator,
	}
}

// GetMe handles GET /api/v1/users/me
// Get current user's profile
func (h *FrontendHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "user profile retrieved successfully", user)
}

// UpdateMe handles PUT /api/v1/users/me
// Update current user's profile
func (h *FrontendHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), userID, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "profile updated successfully", user)
}
