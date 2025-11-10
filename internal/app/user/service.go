package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries      *db.Queries
	auditService *audit.Service
}

func NewService(queries *db.Queries, auditService *audit.Service) *Service {
	return &Service{
		queries:      queries,
		auditService: auditService,
	}
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, errors.Internal("failed to hash password", err)
	}

	// Check if user with email already exists
	_, err = s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.AlreadyExists("user with this email already exists")
	} else if err != pgx.ErrNoRows {
		slog.Error("failed to check existing user by email", "error", err)
		return nil, errors.Internal("failed to check existing user", err)
	}

	// Check if user with username already exists
	_, err = s.queries.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.AlreadyExists("user with this username already exists")
	} else if err != pgx.ErrNoRows {
		slog.Error("failed to check existing user by username", "error", err)
		return nil, errors.Internal("failed to check existing user", err)
	}

	// Create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    pgtype.Text{String: req.FirstName, Valid: req.FirstName != ""},
		LastName:     pgtype.Text{String: req.LastName, Valid: req.LastName != ""},
	})
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, errors.Internal("failed to create user", err)
	}

	// Audit log the user creation
	userID := uuid.UUID(user.ID.Bytes)
	s.auditService.LogCreate(ctx, "users", userID, user)

	return toUserResponse(&user), nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	user, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("user not found")
	} else if err != nil {
		slog.Error("failed to get user", "id", id, "error", err)
		return nil, errors.Internal("failed to get user", err)
	}

	return toUserResponse(&user), nil
}

// ListUsers retrieves a list of users with pagination
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]*UserResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		slog.Error("failed to list users", "error", err)
		return nil, errors.Internal("failed to list users", err)
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		u := user
		responses[i] = toUserResponse(&u)
	}

	return responses, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	// Check if user exists
	oldUser, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == pgx.ErrNoRows {
		return nil, errors.NotFound("user not found")
	} else if err != nil {
		slog.Error("failed to get user", "id", id, "error", err)
		return nil, errors.Internal("failed to get user", err)
	}

	// Update user
	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:        pgtype.UUID{Bytes: userID, Valid: true},
		Email:     pgtype.Text{String: ptrToString(req.Email), Valid: req.Email != nil},
		Username:  pgtype.Text{String: ptrToString(req.Username), Valid: req.Username != nil},
		FirstName: pgtype.Text{String: ptrToString(req.FirstName), Valid: req.FirstName != nil},
		LastName:  pgtype.Text{String: ptrToString(req.LastName), Valid: req.LastName != nil},
	})
	if err != nil {
		slog.Error("failed to update user", "id", id, "error", err)
		return nil, errors.Internal("failed to update user", err)
	}

	// Audit log the user update
	s.auditService.LogUpdate(ctx, "users", userID, oldUser, user)

	return toUserResponse(&user), nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.Validation("invalid user ID format")
	}

	// Check if user exists
	user, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == pgx.ErrNoRows {
		return errors.NotFound("user not found")
	} else if err != nil {
		slog.Error("failed to get user", "id", id, "error", err)
		return errors.Internal("failed to get user", err)
	}

	err = s.queries.DeleteUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		slog.Error("failed to delete user", "id", id, "error", err)
		return errors.Internal("failed to delete user", err)
	}

	// Audit log the user deletion
	s.auditService.LogDelete(ctx, "users", userID, user)

	return nil
}

// Helper functions
func toUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		ID:        uuid.UUID(user.ID.Bytes).String(),
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
