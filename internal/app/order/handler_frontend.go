package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/user/coc/internal/ctxkeys"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// FrontendHandler handles frontend/customer-facing order operations
// Users can only access and manage their OWN orders
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

// CreateOrder handles POST /api/v1/orders
// Create order for the current authenticated user
func (h *FrontendHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Force the user_id to be the authenticated user's ID
	req.UserID = userID

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "order created successfully", order)
}

// GetOrder handles GET /api/v1/orders/{id}
// Get a specific order (only if it belongs to current user)
func (h *FrontendHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Verify the order belongs to the authenticated user
	if order.UserID != userID {
		response.Error(w, http.StatusForbidden, "you don't have permission to access this order")
		return
	}

	response.JSON(w, http.StatusOK, "order retrieved successfully", order)
}

// ListOrders handles GET /api/v1/orders
// List all orders belonging to the current authenticated user
func (h *FrontendHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)

	if limit <= 0 {
		limit = 10
	}

	// Always filter by current user's ID
	orders, err := h.service.ListOrdersByUserID(r.Context(), userID, int32(limit), int32(offset))
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "orders retrieved successfully", orders)
}

// UpdateOrder handles PUT /api/v1/orders/{id}
// Update an order (only if it belongs to current user)
func (h *FrontendHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := ctxkeys.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order ID is required")
		return
	}

	// First verify the order belongs to the user
	existingOrder, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	if existingOrder.UserID != userID {
		response.Error(w, http.StatusForbidden, "you don't have permission to update this order")
		return
	}

	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	order, err := h.service.UpdateOrder(r.Context(), id, req)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "order updated successfully", order)
}
