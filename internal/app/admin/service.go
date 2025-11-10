package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/app/admin_auth"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

// Service handles admin management operations (CRUD)
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

// CreateAdmin creates a new admin (only super_admin should be able to do this)
func (s *Service) CreateAdmin(ctx context.Context, req CreateAdminRequest) (*admin_auth.AdminResponse, error) {
	// Check if admin with email already exists
	_, err := s.queries.GetAdminByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.AlreadyExists("admin with this email already exists")
	}

	// Check if admin with username already exists
	_, err = s.queries.GetAdminByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.AlreadyExists("admin with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Internal("failed to hash password", err)
	}

	// Create admin
	params := db.CreateAdminParams{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	if req.FirstName != "" {
		params.FirstName = pgtype.Text{String: req.FirstName, Valid: true}
	}

	if req.LastName != "" {
		params.LastName = pgtype.Text{String: req.LastName, Valid: true}
	}

	admin, err := s.queries.CreateAdmin(ctx, params)
	if err != nil {
		return nil, errors.Internal("failed to create admin", err)
	}

	// Audit log
	adminID := uuid.UUID(admin.ID.Bytes)
	s.auditService.LogCreate(ctx, "admins", adminID, admin)

	return toAdminResponse(&admin), nil
}

// GetAdmin retrieves an admin by ID
func (s *Service) GetAdmin(ctx context.Context, id string) (*admin_auth.AdminResponse, error) {
	adminUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Internal("invalid admin ID format", err)
	}

	var adminID pgtype.UUID
	adminID.Bytes = adminUUID
	adminID.Valid = true

	admin, err := s.queries.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, errors.NotFound("admin not found")
	}

	return toAdminResponse(&admin), nil
}

// ListAdmins retrieves all admins with pagination
func (s *Service) ListAdmins(ctx context.Context, limit, offset int32) ([]*admin_auth.AdminResponse, error) {
	admins, err := s.queries.ListAdmins(ctx, db.ListAdminsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.Internal("failed to list admins", err)
	}

	responses := make([]*admin_auth.AdminResponse, len(admins))
	for i, admin := range admins {
		responses[i] = toAdminResponse(&admin)
	}

	return responses, nil
}

// UpdateAdmin updates an admin
func (s *Service) UpdateAdmin(ctx context.Context, id string, req UpdateAdminRequest) (*admin_auth.AdminResponse, error) {
	adminUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Internal("invalid admin ID format", err)
	}

	var adminID pgtype.UUID
	adminID.Bytes = adminUUID
	adminID.Valid = true

	// Get old data for audit
	oldAdmin, err := s.queries.GetAdminByID(ctx, adminID)
	if err != nil {
		return nil, errors.NotFound("admin not found")
	}

	// Prepare update params (SQLC UpdateAdminParams uses concrete types)
	params := db.UpdateAdminParams{
		ID:           adminID,
		Email:        oldAdmin.Email,        // default to existing
		Username:     oldAdmin.Username,     // default to existing
		PasswordHash: oldAdmin.PasswordHash, // default to existing
		FirstName:    oldAdmin.FirstName,    // default to existing
		LastName:     oldAdmin.LastName,     // default to existing
		Role:         oldAdmin.Role,         // default to existing
		IsActive:     oldAdmin.IsActive,     // default to existing
	}

	// Update email if provided
	if req.Email != "" {
		params.Email = req.Email
	}

	// Update username if provided
	if req.Username != "" {
		params.Username = req.Username
	}

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.Internal("failed to hash password", err)
		}
		params.PasswordHash = string(hashedPassword)
	}

	// Update first name
	if req.FirstName != "" {
		params.FirstName = pgtype.Text{String: req.FirstName, Valid: true}
	}

	// Update last name
	if req.LastName != "" {
		params.LastName = pgtype.Text{String: req.LastName, Valid: true}
	}

	// Update role if provided
	if req.Role != "" {
		params.Role = req.Role
	}

	// Update is_active if provided
	if req.IsActive != nil {
		params.IsActive = *req.IsActive
	}

	admin, err := s.queries.UpdateAdmin(ctx, params)
	if err != nil {
		return nil, errors.Internal("failed to update admin", err)
	}

	// Audit log
	s.auditService.LogUpdate(ctx, "admins", adminUUID, oldAdmin, admin)

	return toAdminResponse(&admin), nil
}

// DeleteAdmin soft-deletes an admin (sets is_active = false)
func (s *Service) DeleteAdmin(ctx context.Context, id string) error {
	adminUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.Internal("invalid admin ID format", err)
	}

	var adminID pgtype.UUID
	adminID.Bytes = adminUUID
	adminID.Valid = true

	// Get old data for audit
	oldAdmin, err := s.queries.GetAdminByID(ctx, adminID)
	if err != nil {
		return errors.NotFound("admin not found")
	}

	// Soft delete (set is_active = false)
	err = s.queries.DeleteAdmin(ctx, adminID)
	if err != nil {
		return errors.Internal("failed to delete admin", err)
	}

	// Audit log
	s.auditService.LogDelete(ctx, "admins", adminUUID, oldAdmin)

	return nil
}

// Helper function to convert db.Admin to AdminResponse
func toAdminResponse(admin *db.Admin) *admin_auth.AdminResponse {
	adminID, _ := uuid.FromBytes(admin.ID.Bytes[:])

	resp := &admin_auth.AdminResponse{
		ID:       adminID.String(),
		Email:    admin.Email,
		Username: admin.Username,
		Role:     admin.Role,
		IsActive: admin.IsActive,
	}

	if admin.FirstName.Valid {
		resp.FirstName = admin.FirstName.String
	}

	if admin.LastName.Valid {
		resp.LastName = admin.LastName.String
	}

	return resp
}
