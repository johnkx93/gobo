package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

// AdminHandler handles admin-specific order operations
// Admins can manage ALL orders (create, read, update, delete any order)
type AdminHandler struct {
	service  *Service
	validate *validation.Validator
}

func NewAdminHandler(service *Service, validator *validation.Validator) *AdminHandler {
	return &AdminHandler{
		service:  service,
		validate: validator,
	}
}

// CreateOrder handles POST /api/admin/v1/orders
// Admin can create order for ANY user
func (h *AdminHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

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

// GetOrder handles GET /api/admin/v1/orders/{id}
// Admin can get ANY order by ID
func (h *AdminHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	response.JSON(w, http.StatusOK, "order retrieved successfully", order)
}

// ListOrders handles GET /api/admin/v1/orders
// Admin can list ALL orders with optional user_id filter and pagination
func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)
	userID := r.URL.Query().Get("user_id")

	if limit <= 0 {
		limit = 10
	}

	var orders []*OrderResponse
	var err error

	if userID != "" {
		// Admin filtering orders by specific user
		orders, err = h.service.ListOrdersByUserID(r.Context(), userID, int32(limit), int32(offset))
	} else {
		// Admin viewing all orders
		orders, err = h.service.ListOrders(r.Context(), int32(limit), int32(offset))
	}

	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "orders retrieved successfully", orders)
}

// UpdateOrder handles PUT /api/admin/v1/orders/{id}
// Admin can update ANY order
func (h *AdminHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order ID is required")
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

// DeleteOrder handles DELETE /api/admin/v1/orders/{id}
// Admin can delete ANY order
func (h *AdminHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "order ID is required")
		return
	}

	err := h.service.DeleteOrder(r.Context(), id)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "order deleted successfully", nil)
}
