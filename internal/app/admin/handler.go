package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// Handler handles admin management operations (CRUD on admins)
type Handler struct {
	service  *Service
	validate *validation.Validator
}

func NewHandler(service *Service, validator *validation.Validator) *Handler {
	return &Handler{
		service:  service,
		validate: validator,
	}
}

// CreateAdmin handles POST /api/admin/v1/admins
// Only super_admin should be able to create new admins
// @Summary      Create admin (admin management)
// @Description  Create a new admin user (requires super_admin role)
// @Tags         Admin Management
// @Accept       json
// @Produce      json
// @Param        request body CreateAdminRequest true "Admin data"
// @Success      201 {object} response.JSONResponse{data=internal_app_admin_auth.AdminResponse} "Admin created successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/admin/v1/admins [post]
func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	admin, err := h.service.CreateAdmin(r.Context(), req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "admin created successfully", admin)
}

// GetAdmin handles GET /api/admin/v1/admins/{id}
// @Summary      Get admin (admin management)
// @Description  Retrieve an admin user by ID
// @Tags         Admin Management
// @Accept       json
// @Produce      json
// @Param        id path string true "Admin ID"
// @Success      200 {object} response.JSONResponse{data=internal_app_admin_auth.AdminResponse} "Admin retrieved successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Failure      404 {object} response.JSONResponse "Admin not found"
// @Security     BearerAuth
// @Router       /api/admin/v1/admins/{id} [get]
func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "admin ID is required")
		return
	}

	admin, err := h.service.GetAdmin(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "admin retrieved successfully", admin)
}

// ListAdmins handles GET /api/admin/v1/admins
// @Summary      List admins (admin management)
// @Description  Retrieve all admin users with pagination
// @Tags         Admin Management
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of admins to return (default 10)"
// @Param        offset query int false "Number of admins to skip (default 0)"
// @Success      200 {object} response.JSONResponse{data=[]internal_app_admin_auth.AdminResponse} "Admins retrieved successfully"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/admin/v1/admins [get]
func (h *Handler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)

	if limit <= 0 {
		limit = 10
	}

	admins, err := h.service.ListAdmins(r.Context(), int32(limit), int32(offset))
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "admins retrieved successfully", admins)
}

// UpdateAdmin handles PUT /api/admin/v1/admins/{id}
// @Summary      Update admin (admin management)
// @Description  Update an admin user by ID
// @Tags         Admin Management
// @Accept       json
// @Produce      json
// @Param        id path string true "Admin ID"
// @Param        request body UpdateAdminRequest true "Admin update data"
// @Success      200 {object} response.JSONResponse{data=internal_app_admin_auth.AdminResponse} "Admin updated successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Failure      404 {object} response.JSONResponse "Admin not found"
// @Security     BearerAuth
// @Router       /api/admin/v1/admins/{id} [put]
func (h *Handler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "admin ID is required")
		return
	}

	var req UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	admin, err := h.service.UpdateAdmin(r.Context(), id, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "admin updated successfully", admin)
}

// DeleteAdmin handles DELETE /api/admin/v1/admins/{id}
// @Summary      Delete admin (admin management)
// @Description  Delete an admin user by ID
// @Tags         Admin Management
// @Accept       json
// @Produce      json
// @Param        id path string true "Admin ID"
// @Success      200 {object} response.JSONResponse "Admin deleted successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Failure      404 {object} response.JSONResponse "Admin not found"
// @Security     BearerAuth
// @Router       /api/admin/v1/admins/{id} [delete]
func (h *Handler) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "admin ID is required")
		return
	}

	err := h.service.DeleteAdmin(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "admin deleted successfully", nil)
}
