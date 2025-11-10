package admin_auth

// LoginRequest represents admin login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents admin login response
type LoginResponse struct {
	Token string         `json:"token"`
	Admin *AdminResponse `json:"admin"`
}

// AdminResponse represents admin data in response
type AdminResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
}
