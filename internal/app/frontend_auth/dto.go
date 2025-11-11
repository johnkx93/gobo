package frontend_auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// UserResponse represents the user data in responses (excluding password)
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts a db.User to UserResponse (excluding password_hash)
func ToUserResponse(id pgtype.UUID, email, username string, firstName, lastName pgtype.Text, createdAt, updatedAt pgtype.Timestamptz) UserResponse {
	resp := UserResponse{
		Email:     email,
		Username:  username,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}

	// Convert UUID to string
	if id.Valid {
		userID, err := uuid.FromBytes(id.Bytes[:])
		if err == nil {
			resp.ID = userID.String()
		}
	}

	if firstName.Valid {
		resp.FirstName = &firstName.String
	}

	if lastName.Valid {
		resp.LastName = &lastName.String
	}

	return resp
}
