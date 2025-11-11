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
	service  *Service
	validate *validation.Validator
}

func NewFrontendHandler(service *Service, validator *validation.Validator) *FrontendHandler {
	return &FrontendHandler{
		service:  service,
		validate: validator,
	}
}

// CreateMyAddress handles POST /api/v1/addresses
// User creates their own address
func (h *FrontendHandler) CreateMyAddress(w http.ResponseWriter, r *http.Request) {
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

	var req CreateAddressForUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	_, err = h.service.CreateAddressForUser(r.Context(), userID, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after creating
	addresses, err := h.service.ListMyAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "address created successfully", addresses)
}

// GetMyAddress handles GET /api/v1/addresses/{id}
// User gets their own address by ID
func (h *FrontendHandler) GetMyAddress(w http.ResponseWriter, r *http.Request) {
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

	address, err := h.service.GetAddressForUser(r.Context(), userID, addressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address retrieved successfully", address)
}

// ListMyAddresses handles GET /api/v1/addresses
// User lists all their own addresses
func (h *FrontendHandler) ListMyAddresses(w http.ResponseWriter, r *http.Request) {
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

	addresses, err := h.service.ListMyAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "addresses retrieved successfully", addresses)
}

// UpdateMyAddress handles PUT /api/v1/addresses/{id}
// User updates their own address
func (h *FrontendHandler) UpdateMyAddress(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.service.UpdateAddressForUser(r.Context(), userID, addressID, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after updating
	addresses, err := h.service.ListMyAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address updated successfully", addresses)
}

// DeleteMyAddress handles DELETE /api/v1/addresses/{id}
// User deletes their own address
func (h *FrontendHandler) DeleteMyAddress(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteAddressForUser(r.Context(), userID, addressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after deleting
	addresses, err := h.service.ListMyAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "address deleted successfully", addresses)
}

// SetMyDefaultAddress handles POST /api/v1/addresses/default
// User sets their default address
func (h *FrontendHandler) SetMyDefaultAddress(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.SetMyDefaultAddress(r.Context(), userID, req.AddressID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Return all addresses for the user after setting default
	addresses, err := h.service.ListMyAddresses(r.Context(), userID)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "default address set successfully", addresses)
}
