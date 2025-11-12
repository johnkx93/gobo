package admin

// CreateAdminRequest represents the request to create a new admin
type CreateAdminRequest struct {
	Email     string `json:"email" validate:"required,email" example:"newadmin@example.com"`
	Username  string `json:"username" validate:"required,min=3,max=50" example:"newadmin"`
	Password  string `json:"password" validate:"required,min=8" example:"SecurePass123"`
	FirstName string `json:"first_name" validate:"required,max=100" example:"John"`
	LastName  string `json:"last_name" validate:"required,max=100" example:"Doe"`
	// Role is optional on request; if omitted the service will default to "moderator".
	Role string `json:"role,omitempty" validate:"omitempty,oneof=admin super_admin moderator" example:"moderator"`
}

// UpdateAdminRequest represents the request to update an admin
type UpdateAdminRequest struct {
	Email     string `json:"email,omitempty" validate:"omitempty,email" example:"updated@example.com"`
	Username  string `json:"username,omitempty" validate:"omitempty,min=3,max=50" example:"updatedadmin"`
	Password  string `json:"password,omitempty" validate:"omitempty,min=8" example:"NewSecurePass123"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=100" example:"Jane"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=100" example:"Smith"`
	Role      string `json:"role,omitempty" validate:"omitempty,oneof=admin super_admin moderator" example:"admin"`
	IsActive  *bool  `json:"is_active,omitempty" example:"true"`
}
