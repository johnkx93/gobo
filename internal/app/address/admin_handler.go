package address

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// AdminHandler handles admin-specific address operations
// Admins can manage ALL addresses (create, read, update, delete any address)
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

// CreateAddress handles POST /api/admin/v1/addresses
// Admin creates a new address for any user
func (h *AdminHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	var req CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	address, err := h.service.CreateAddress(r.Context(), req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "address created successfully", address)
}

// GetAddress handles GET /api/admin/v1/addresses/{id}
// Admin can get ANY address by ID
func (h *AdminHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	address, err := h.service.GetAddress(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address retrieved successfully", address)
}

// ListAllAddresses handles GET /api/admin/v1/addresses
// Admin can list ALL addresses with pagination
func (h *AdminHandler) ListAllAddresses(w http.ResponseWriter, r *http.Request) {
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

	addresses, err := h.service.ListAllAddresses(r.Context(), int32(limit), int32(offset))
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "addresses retrieved successfully", addresses)
}

// ListAddressesByUser handles GET /api/admin/v1/users/{user_id}/addresses
// Admin can list all addresses for a specific user
func (h *AdminHandler) ListAddressesByUser(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	addresses, err := h.service.ListAddressesByUser(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "addresses retrieved successfully", addresses)
}

// UpdateAddress handles PUT /api/admin/v1/addresses/{id}
// Admin can update ANY address
func (h *AdminHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	var req UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	address, err := h.service.UpdateAddress(r.Context(), id, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address updated successfully", address)
}

// DeleteAddress handles DELETE /api/admin/v1/addresses/{id}
// Admin can delete ANY address
func (h *AdminHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	err := h.service.DeleteAddress(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address deleted successfully", nil)
}

// SetDefaultAddress handles POST /api/admin/v1/users/{user_id}/addresses/default
// Admin can set default address for any user
func (h *AdminHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	// REQUIRED: Check admin role first
	role, ok := ctxkeys.GetAdminRole(r)
	if !ok || role == "" {
		response.Error(w, http.StatusUnauthorized, "admin role not found")
		return
	}

	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	var req SetDefaultAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	err := h.service.SetDefaultAddress(r.Context(), userID, req.AddressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "default address set successfully", nil)
}
