package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/coc/pkg/response"
	"github.com/user/coc/pkg/validation"
)

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

// CreateOrder handles POST /orders
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
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

// GetOrder handles GET /orders/{id}
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

// ListOrders handles GET /orders
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)
	userID := r.URL.Query().Get("user_id")

	if limit <= 0 {
		limit = 10
	}

	var orders []*OrderResponse
	var err error

	if userID != "" {
		orders, err = h.service.ListOrdersByUserID(r.Context(), userID, int32(limit), int32(offset))
	} else {
		orders, err = h.service.ListOrders(r.Context(), int32(limit), int32(offset))
	}

	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "orders retrieved successfully", orders)
}

// UpdateOrder handles PUT /orders/{id}
func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
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

// DeleteOrder handles DELETE /orders/{id}
func (h *Handler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
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
