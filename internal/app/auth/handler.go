package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
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

// Login handles POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	token, user, err := h.service.Login(r.Context(), strings.TrimSpace(req.Email), strings.TrimSpace(req.Password))
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Convert user to response (excluding password)
	userResp := ToUserResponse(
		user.ID,
		user.Email,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	loginResp := LoginResponse{
		Token: token,
		User:  userResp,
	}

	response.JSON(w, http.StatusOK, "login successful", loginResp)
}

// Register handles POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	// decode request body into DTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errorMsg := h.validate.TranslateErrors(err)
		response.Error(w, http.StatusBadRequest, errorMsg)
		return
	}

	user, err := h.service.Register(
		r.Context(),
		strings.TrimSpace(req.Email),
		strings.TrimSpace(req.Username),
		strings.TrimSpace(req.Password),
		strings.TrimSpace(req.FirstName),
		strings.TrimSpace(req.LastName),
	)

	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Generate token for the new user
	token, err := h.service.GenerateToken(user)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Convert user to response (excluding password)
	userResp := ToUserResponse(
		user.ID,
		user.Email,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	registerResp := LoginResponse{
		Token: token,
		User:  userResp,
	}

	response.JSON(w, http.StatusCreated, "registration successful", registerResp)
}

// UserFromContext extracts user from request context (set by auth middleware)
func UserFromContext(r *http.Request) *db.User {
	user, ok := r.Context().Value("user").(*db.User)
	if !ok {
		return nil
	}
	return user
}
