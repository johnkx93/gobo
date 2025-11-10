package admin_management

// CreateAdminRequest represents the request to create a new admin
type CreateAdminRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	// Role is optional on request; if omitted the service will default to "moderator".
	Role string `json:"role,omitempty" validate:"omitempty,oneof=admin super_admin moderator"`
}

// UpdateAdminRequest represents the request to update an admin
type UpdateAdminRequest struct {
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	Username  string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Password  string `json:"password,omitempty" validate:"omitempty,min=8"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	Role      string `json:"role,omitempty" validate:"omitempty,oneof=admin super_admin moderator"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
