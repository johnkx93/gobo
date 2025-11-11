package user

import (
	"context"

	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
)

// FrontendService contains business logic for frontend (user) operations
// Frontend users can only operate on their own entities
type FrontendService struct {
	adminService *AdminService
	queries      *db.Queries
	auditService *audit.Service
}

func NewFrontendService(queries *db.Queries, auditService *audit.Service) *FrontendService {
	return &FrontendService{
		adminService: NewAdminService(queries, auditService),
		queries:      queries,
		auditService: auditService,
	}
}

// GetUser returns the user's own profile
func (s *FrontendService) GetUser(ctx context.Context, userID string) (*UserResponse, error) {
	// For frontend, userID is the only ID we operate on
	return s.adminService.GetUser(ctx, userID)
}

// UpdateUser updates the current user's profile
func (s *FrontendService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*UserResponse, error) {
	return s.adminService.UpdateUser(ctx, userID, req)
}
