package address

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// FrontendHandler handles frontend/customer-facing address operations
// Users can only access and modify their OWN addresses
type FrontendHandler struct {
	service  *UserService
	validate *validation.Validator
}

func NewFrontendHandler(service *UserService, validator *validation.Validator) *FrontendHandler {
	return &FrontendHandler{
		service:  service,
		validate: validator,
	}
}

// CreateAddress handles POST /api/v1/addresses
// User creates their own address
func (h *FrontendHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var req UserCreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	_, err = h.service.CreateAddress(r.Context(), userID, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after creating
	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "address created successfully", addresses)
}

// GetAddress handles GET /api/v1/addresses/{id}
// User gets their own address by ID
func (h *FrontendHandler) GetAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	addressID := chi.URLParam(r, "id")
	if addressID == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	address, err := h.service.GetAddress(r.Context(), userID, addressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address retrieved successfully", address)
}

// ListAddresses handles GET /api/v1/addresses
// User lists all their own addresses
func (h *FrontendHandler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "addresses retrieved successfully", addresses)
}

// UpdateAddress handles PUT /api/v1/addresses/{id}
// User updates their own address
func (h *FrontendHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	addressID := chi.URLParam(r, "id")
	if addressID == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	var req UserUpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	_, err = h.service.UpdateAddress(r.Context(), userID, addressID, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after updating
	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address updated successfully", addresses)
}

// DeleteAddress handles DELETE /api/v1/addresses/{id}
// User deletes their own address
func (h *FrontendHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	addressID := chi.URLParam(r, "id")
	if addressID == "" {
		response.Error(w, http.StatusBadRequest, "address ID is required")
		return
	}

	err = h.service.DeleteAddress(r.Context(), userID, addressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after deleting
	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address deleted successfully", addresses)
}

// SetDefaultAddress handles POST /api/v1/addresses/default
// User sets their default address
func (h *FrontendHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
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

	err = h.service.SetDefaultAddress(r.Context(), userID, req.AddressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after setting default
	addresses, err := h.service.ListAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "default address set successfully", addresses)
}
