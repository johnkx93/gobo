package user

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=100"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"omitempty,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,max=100"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Username  *string `json:"username" validate:"omitempty,min=3,max=100"`
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListUsersRequest represents the request to list users
type ListUsersRequest struct {
	Limit  int32 `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int32 `json:"offset" validate:"omitempty,min=0"`
}
