package admin_auth

// LoginRequest represents admin login request
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"admin"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse represents admin login response
type LoginResponse struct {
	Token string         `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Admin *AdminResponse `json:"admin"`
}

// AdminResponse represents admin data in response
type AdminResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"admin@example.com"`
	Username  string `json:"username" example:"admin"`
	FirstName string `json:"first_name,omitempty" example:"John"`
	LastName  string `json:"last_name,omitempty" example:"Doe"`
	Role      string `json:"role" example:"super_admin"`
	IsActive  bool   `json:"is_active" example:"true"`
}
