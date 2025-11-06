package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/user/coc/internal/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

// JSONResponse represents a standard JSON response
type JSONResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, "user created successfully", user)
}

// GetUser handles GET /users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, "user retrieved successfully", user)
}

// ListUsers handles GET /users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)

	if limit <= 0 {
		limit = 10
	}

	users, err := h.service.ListUsers(r.Context(), int32(limit), int32(offset))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, "users retrieved successfully", users)
}

// UpdateUser handles PUT /users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, "user updated successfully", user)
}

// DeleteUser handles DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, "user deleted successfully", nil)
}

// Helper methods
func (h *Handler) respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  false,
		Message: message,
	})
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	var domainErr *errors.DomainError
	if e, ok := err.(*errors.DomainError); ok {
		domainErr = e
	} else {
		slog.Error("unexpected error type", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch domainErr.Code {
	case errors.CodeNotFound:
		h.respondError(w, http.StatusNotFound, domainErr.Message)
	case errors.CodeAlreadyExists:
		h.respondError(w, http.StatusConflict, domainErr.Message)
	case errors.CodeValidation:
		h.respondError(w, http.StatusBadRequest, domainErr.Message)
	case errors.CodeUnauthorized:
		h.respondError(w, http.StatusUnauthorized, domainErr.Message)
	default:
		slog.Error("internal error", "error", domainErr)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
