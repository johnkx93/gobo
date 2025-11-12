package user

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Username  string `json:"username" validate:"required,min=3,max=100" example:"johndoe"`
	Password  string `json:"password" validate:"required,min=8" example:"SecurePass123"`
	FirstName string `json:"first_name" validate:"omitempty,max=100" example:"John"`
	LastName  string `json:"last_name" validate:"omitempty,max=100" example:"Doe"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email" example:"newemail@example.com"`
	Username  *string `json:"username" validate:"omitempty,min=3,max=100" example:"newusername"`
	FirstName *string `json:"first_name" validate:"omitempty,max=100" example:"Jane"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100" example:"Smith"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID               string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email            string `json:"email" example:"john.doe@example.com"`
	Username         string `json:"username" example:"johndoe"`
	FirstName        string `json:"first_name,omitempty" example:"John"`
	LastName         string `json:"last_name,omitempty" example:"Doe"`
	DefaultAddressID string `json:"default_address_id,omitempty" example:"650e8400-e29b-41d4-a716-446655440001"`
	CreatedAt        string `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt        string `json:"updated_at" example:"2024-01-02T15:30:00Z"`
}

// ListUsersRequest represents the request to list users
type ListUsersRequest struct {
	Limit  int32 `json:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset int32 `json:"offset" validate:"omitempty,min=0" example:"0"`
}
