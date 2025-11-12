package user

import (
	"encoding/json"
	"net/http"

	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// FrontendHandler handles frontend/customer-facing user operations
// Users can only access and modify their OWN data
type FrontendHandler struct {
	service  *FrontendService
	validate *validation.Validator
}

func NewFrontendHandler(service *FrontendService, validator *validation.Validator) *FrontendHandler {
	return &FrontendHandler{
		service:  service,
		validate: validator,
	}
}

// GetMe returns the current user's profile
// @Summary      Get current user profile
// @Description  Retrieve the authenticated user's profile information
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} response.JSONResponse{data=UserResponse} "User profile retrieved"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/users/me [get]
func (h *FrontendHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
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
// @Summary      Update current user profile
// @Description  Update the authenticated user's profile information
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Param        request body UpdateUserRequest true "Profile update data"
// @Success      200 {object} response.JSONResponse{data=UserResponse} "Profile updated successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/users/me [put]
func (h *FrontendHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
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
