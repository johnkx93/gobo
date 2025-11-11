package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// AdminHandler handles admin-specific user operations
// Admins can manage ALL users (create, read, update, delete any user)
type AdminHandler struct {
	service  *AdminService
	validate *validation.Validator
}

func NewAdminHandler(service *AdminService, validator *validation.Validator) *AdminHandler {
	return &AdminHandler{
		service:  service,
		validate: validator,
	}
}

// CreateUser handles POST /api/admin/v1/users
// Admin creates a new user (can create any user)
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "user created successfully", user)
}

// GetUser handles GET /api/admin/v1/users/{id}
// Admin can get ANY user by ID
func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "user retrieved successfully", user)
}

// ListUsers handles GET /api/admin/v1/users
// Admin can list ALL users with pagination
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)

	if limit <= 0 {
		limit = 10
	}

	users, err := h.service.ListUsers(r.Context(), int32(limit), int32(offset))
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "users retrieved successfully", users)
}

// UpdateUser handles PUT /api/admin/v1/users/{id}
// Admin can update ANY user
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
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

	user, err := h.service.UpdateUser(r.Context(), id, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "user updated successfully", user)
}

// DeleteUser handles DELETE /api/admin/v1/users/{id}
// Admin can delete ANY user
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "user deleted successfully", nil)
}
