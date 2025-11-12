package admin_auth

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/user/coc/internal/response"
	"github.com/user/coc/internal/validation"
)

type AuthHandler struct {
	authService *AuthService
	validate    *validation.Validator
}

func NewAuthHandler(authService *AuthService, validator *validation.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator,
	}
}

// Login handles POST /api/admin/v1/auth/login
// @Summary      Admin login
// @Description  Authenticate admin user and receive JWT token
// @Tags         Admin Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} response.JSONResponse{data=LoginResponse} "Login successful"
// @Failure      400 {object} response.JSONResponse "Invalid request or credentials"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Router       /api/admin/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	token, admin, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		response.HandleServiceError(w, err)
		return
	}

	// Convert UUID to string
	adminID, _ := uuid.FromBytes(admin.ID.Bytes[:])

	adminResp := &AdminResponse{
		ID:       adminID.String(),
		Email:    admin.Email,
		Username: admin.Username,
		Role:     admin.Role,
		IsActive: admin.IsActive,
	}

	if admin.FirstName.Valid {
		adminResp.FirstName = admin.FirstName.String
	}
	if admin.LastName.Valid {
		adminResp.LastName = admin.LastName.String
	}

	resp := &LoginResponse{
		Token: token,
		Admin: adminResp,
	}

	response.JSON(w, http.StatusOK, "admin login successful", resp)
}
