package address

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
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

// CreateAddress creates a new address (admin only - can create for any user)
func (s *Service) CreateAddress(ctx context.Context, req CreateAddressRequest) (*AddressResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	// Verify user exists
	_, err = s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}
		slog.Error("failed to verify user exists", "user_id", req.UserID, "error", err)
		return nil, errors.Internal("failed to verify user", err)
	}

	// Create address
	address, err := s.queries.CreateAddress(ctx, db.CreateAddressParams{
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		Address:     req.Address,
		Floor:       req.Floor,
		UnitNo:      req.UnitNo,
		BlockTower:  pgtype.Text{String: getStringValue(req.BlockTower), Valid: req.BlockTower != nil},
		CompanyName: pgtype.Text{String: getStringValue(req.CompanyName), Valid: req.CompanyName != nil},
	})
	if err != nil {
		slog.Error("failed to create address", "error", err)
		return nil, errors.Internal("failed to create address", err)
	}

	// Audit log the address creation
	addressID := uuid.UUID(address.ID.Bytes)
	s.auditService.LogCreate(ctx, "addresses", addressID, address)

	return toAddressResponse(&address), nil
}

// CreateAddressForUser creates a new address for the authenticated user
func (s *Service) CreateAddressForUser(ctx context.Context, userID uuid.UUID, req CreateAddressForUserRequest) (*AddressResponse, error) {
	// Create address
	address, err := s.queries.CreateAddress(ctx, db.CreateAddressParams{
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		Address:     req.Address,
		Floor:       req.Floor,
		UnitNo:      req.UnitNo,
		BlockTower:  pgtype.Text{String: getStringValue(req.BlockTower), Valid: req.BlockTower != nil},
		CompanyName: pgtype.Text{String: getStringValue(req.CompanyName), Valid: req.CompanyName != nil},
	})
	if err != nil {
		slog.Error("failed to create address for user", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to create address", err)
	}

	// Audit log the address creation
	addressID := uuid.UUID(address.ID.Bytes)
	s.auditService.LogCreate(ctx, "addresses", addressID, address)

	return toAddressResponse(&address), nil
}

// GetAddress retrieves an address by ID (admin only)
func (s *Service) GetAddress(ctx context.Context, id string) (*AddressResponse, error) {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid address ID format")
	}

	address, err := s.queries.GetAddressByID(ctx, pgtype.UUID{Bytes: addressID, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("address not found")
		}
		slog.Error("failed to get address", "address_id", id, "error", err)
		return nil, errors.Internal("failed to get address", err)
	}

	return toAddressResponse(&address), nil
}

// GetAddressForUser retrieves an address by ID for the authenticated user (user can only get their own)
func (s *Service) GetAddressForUser(ctx context.Context, userID uuid.UUID, addressID string) (*AddressResponse, error) {
	addrID, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.Validation("invalid address ID format")
	}

	address, err := s.queries.GetAddressByIDAndUserID(ctx, db.GetAddressByIDAndUserIDParams{
		ID:     pgtype.UUID{Bytes: addrID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("address not found")
		}
		slog.Error("failed to get address for user", "user_id", userID, "address_id", addressID, "error", err)
		return nil, errors.Internal("failed to get address", err)
	}

	return toAddressResponse(&address), nil
}

// ListAddressesByUser retrieves all addresses for a specific user
func (s *Service) ListAddressesByUser(ctx context.Context, userID string) ([]*AddressResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.Validation("invalid user ID format")
	}

	addresses, err := s.queries.GetAddressesByUserID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		slog.Error("failed to list addresses by user", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to list addresses", err)
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, addr := range addresses {
		a := addr
		responses[i] = toAddressResponse(&a)
	}

	return responses, nil
}

// ListMyAddresses retrieves all addresses for the authenticated user
func (s *Service) ListMyAddresses(ctx context.Context, userID uuid.UUID) ([]*AddressResponse, error) {
	addresses, err := s.queries.GetAddressesByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		slog.Error("failed to list my addresses", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to list addresses", err)
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, addr := range addresses {
		a := addr
		responses[i] = toAddressResponse(&a)
	}

	return responses, nil
}

// ListAllAddresses retrieves all addresses with pagination (admin only)
func (s *Service) ListAllAddresses(ctx context.Context, limit, offset int32) ([]*AddressResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	addresses, err := s.queries.ListAllAddresses(ctx, db.ListAllAddressesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		slog.Error("failed to list all addresses", "error", err)
		return nil, errors.Internal("failed to list addresses", err)
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, addr := range addresses {
		a := addr
		responses[i] = toAddressResponse(&a)
	}

	return responses, nil
}

// UpdateAddress updates an address (admin only)
func (s *Service) UpdateAddress(ctx context.Context, id string, req UpdateAddressRequest) (*AddressResponse, error) {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Validation("invalid address ID format")
	}

	// Get old address for audit
	oldAddress, err := s.queries.GetAddressByID(ctx, pgtype.UUID{Bytes: addressID, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("address not found")
		}
		slog.Error("failed to get address for update", "address_id", id, "error", err)
		return nil, errors.Internal("failed to get address", err)
	}

	// Update address
	address, err := s.queries.UpdateAddress(ctx, db.UpdateAddressParams{
		ID:          pgtype.UUID{Bytes: addressID, Valid: true},
		Address:     pgtype.Text{String: getStringValue(req.Address), Valid: req.Address != nil},
		Floor:       pgtype.Text{String: getStringValue(req.Floor), Valid: req.Floor != nil},
		UnitNo:      pgtype.Text{String: getStringValue(req.UnitNo), Valid: req.UnitNo != nil},
		BlockTower:  pgtype.Text{String: getStringValue(req.BlockTower), Valid: req.BlockTower != nil},
		CompanyName: pgtype.Text{String: getStringValue(req.CompanyName), Valid: req.CompanyName != nil},
	})
	if err != nil {
		slog.Error("failed to update address", "address_id", id, "error", err)
		return nil, errors.Internal("failed to update address", err)
	}

	// Audit log the update
	s.auditService.LogUpdate(ctx, "addresses", addressID, oldAddress, address)

	return toAddressResponse(&address), nil
}

// UpdateAddressForUser updates an address for the authenticated user
func (s *Service) UpdateAddressForUser(ctx context.Context, userID uuid.UUID, addressID string, req UpdateAddressRequest) (*AddressResponse, error) {
	addrID, err := uuid.Parse(addressID)
	if err != nil {
		return nil, errors.Validation("invalid address ID format")
	}

	// Get old address for audit (also verifies ownership)
	oldAddress, err := s.queries.GetAddressByIDAndUserID(ctx, db.GetAddressByIDAndUserIDParams{
		ID:     pgtype.UUID{Bytes: addrID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("address not found")
		}
		slog.Error("failed to get address for user update", "user_id", userID, "address_id", addressID, "error", err)
		return nil, errors.Internal("failed to get address", err)
	}

	// Update address
	address, err := s.queries.UpdateAddressForUser(ctx, db.UpdateAddressForUserParams{
		ID:          pgtype.UUID{Bytes: addrID, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		Address:     pgtype.Text{String: getStringValue(req.Address), Valid: req.Address != nil},
		Floor:       pgtype.Text{String: getStringValue(req.Floor), Valid: req.Floor != nil},
		UnitNo:      pgtype.Text{String: getStringValue(req.UnitNo), Valid: req.UnitNo != nil},
		BlockTower:  pgtype.Text{String: getStringValue(req.BlockTower), Valid: req.BlockTower != nil},
		CompanyName: pgtype.Text{String: getStringValue(req.CompanyName), Valid: req.CompanyName != nil},
	})
	if err != nil {
		slog.Error("failed to update address for user", "user_id", userID, "address_id", addressID, "error", err)
		return nil, errors.Internal("failed to update address", err)
	}

	// Audit log the update
	s.auditService.LogUpdate(ctx, "addresses", addrID, oldAddress, address)

	return toAddressResponse(&address), nil
}

// DeleteAddress deletes an address (admin only)
func (s *Service) DeleteAddress(ctx context.Context, id string) error {
	addressID, err := uuid.Parse(id)
	if err != nil {
		return errors.Validation("invalid address ID format")
	}

	// Get address for audit before deletion
	address, err := s.queries.GetAddressByID(ctx, pgtype.UUID{Bytes: addressID, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("address not found")
		}
		slog.Error("failed to get address for deletion", "address_id", id, "error", err)
		return errors.Internal("failed to get address", err)
	}

	// Delete address
	err = s.queries.DeleteAddress(ctx, pgtype.UUID{Bytes: addressID, Valid: true})
	if err != nil {
		slog.Error("failed to delete address", "address_id", id, "error", err)
		return errors.Internal("failed to delete address", err)
	}

	// Audit log the deletion
	s.auditService.LogDelete(ctx, "addresses", addressID, address)

	return nil
}

// DeleteAddressForUser deletes an address for the authenticated user
func (s *Service) DeleteAddressForUser(ctx context.Context, userID uuid.UUID, addressID string) error {
	addrID, err := uuid.Parse(addressID)
	if err != nil {
		return errors.Validation("invalid address ID format")
	}

	// Get address for audit before deletion (also verifies ownership)
	address, err := s.queries.GetAddressByIDAndUserID(ctx, db.GetAddressByIDAndUserIDParams{
		ID:     pgtype.UUID{Bytes: addrID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("address not found")
		}
		slog.Error("failed to get address for user deletion", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to get address", err)
	}

	// Delete address
	err = s.queries.DeleteAddressForUser(ctx, db.DeleteAddressForUserParams{
		ID:     pgtype.UUID{Bytes: addrID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		slog.Error("failed to delete address for user", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to delete address", err)
	}

	// Audit log the deletion
	s.auditService.LogDelete(ctx, "addresses", addrID, address)

	return nil
}

// SetDefaultAddress sets the default address for a user (admin only)
func (s *Service) SetDefaultAddress(ctx context.Context, userID, addressID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.Validation("invalid user ID format")
	}

	addrID, err := uuid.Parse(addressID)
	if err != nil {
		return errors.Validation("invalid address ID format")
	}

	// Verify address belongs to user
	_, err = s.queries.GetAddressByIDAndUserID(ctx, db.GetAddressByIDAndUserIDParams{
		ID:     pgtype.UUID{Bytes: addrID, Valid: true},
		UserID: pgtype.UUID{Bytes: uid, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("address not found or does not belong to user")
		}
		slog.Error("failed to verify address ownership", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to verify address", err)
	}

	// Set default address
	_, err = s.queries.SetDefaultAddress(ctx, db.SetDefaultAddressParams{
		ID:               pgtype.UUID{Bytes: uid, Valid: true},
		DefaultAddressID: pgtype.UUID{Bytes: addrID, Valid: true},
	})
	if err != nil {
		slog.Error("failed to set default address", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to set default address", err)
	}

	return nil
}

// SetMyDefaultAddress sets the default address for the authenticated user
func (s *Service) SetMyDefaultAddress(ctx context.Context, userID uuid.UUID, addressID string) error {
	addrID, err := uuid.Parse(addressID)
	if err != nil {
		return errors.Validation("invalid address ID format")
	}

	// Set default address (query verifies address belongs to user)
	_, err = s.queries.SetDefaultAddressForUser(ctx, db.SetDefaultAddressForUserParams{
		ID:               pgtype.UUID{Bytes: userID, Valid: true},
		DefaultAddressID: pgtype.UUID{Bytes: addrID, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.NotFound("address not found or does not belong to you")
		}
		slog.Error("failed to set my default address", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to set default address", err)
	}

	return nil
}

// Helper functions

func toAddressResponse(address *db.Address) *AddressResponse {
	return &AddressResponse{
		ID:          uuid.UUID(address.ID.Bytes).String(),
		UserID:      uuid.UUID(address.UserID.Bytes).String(),
		Address:     address.Address,
		Floor:       address.Floor,
		UnitNo:      address.UnitNo,
		BlockTower:  getStringPointer(address.BlockTower),
		CompanyName: getStringPointer(address.CompanyName),
		CreatedAt:   address.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   address.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getStringPointer(text pgtype.Text) *string {
	if !text.Valid {
		return nil
	}
	return &text.String
}
