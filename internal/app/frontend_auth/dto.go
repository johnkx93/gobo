package frontend_auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Username  string `json:"username" validate:"required,min=3,max=50" example:"johndoe"`
	Password  string `json:"password" validate:"required,min=6" example:"SecurePass123"`
	FirstName string `json:"first_name,omitempty" example:"John"`
	LastName  string `json:"last_name,omitempty" example:"Doe"`
}

// UserResponse represents the user data in responses (excluding password)
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Username  string    `json:"username" example:"johndoe"`
	FirstName *string   `json:"first_name,omitempty" example:"John"`
	LastName  *string   `json:"last_name,omitempty" example:"Doe"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-02T15:30:00Z"`
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
