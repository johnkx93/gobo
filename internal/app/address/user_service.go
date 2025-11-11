package address

import (
	"context"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
)

// UserService handles user address operations
// Users can only manage their OWN addresses
type UserService struct {
	queries      *db.Queries
	auditService *audit.Service
}

func NewUserService(queries *db.Queries, auditService *audit.Service) *UserService {
	return &UserService{
		queries:      queries,
		auditService: auditService,
	}
}

// CreateAddress creates a new address for the authenticated user
func (s *UserService) CreateAddress(ctx context.Context, userID uuid.UUID, req UserCreateAddressRequest) (*AddressResponse, error) {
	// Create address
	address, err := s.queries.CreateAddress(ctx, db.CreateAddressParams{
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		Address:     strings.TrimSpace(req.Address),
		Floor:       strings.TrimSpace(req.Floor),
		UnitNo:      strings.TrimSpace(req.UnitNo),
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

// GetAddress retrieves an address by ID (only if it belongs to the user)
func (s *UserService) GetAddress(ctx context.Context, userID uuid.UUID, addressID string) (*AddressResponse, error) {
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

// ListAddresses retrieves all addresses for the authenticated user
func (s *UserService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]*AddressResponse, error) {
	addresses, err := s.queries.GetAddressesByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		slog.Error("failed to list user addresses", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to list addresses", err)
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, addr := range addresses {
		a := addr
		responses[i] = toAddressResponse(&a)
	}

	return responses, nil
}

// UpdateAddress updates an address (only if it belongs to the user)
func (s *UserService) UpdateAddress(ctx context.Context, userID uuid.UUID, addressID string, req UserUpdateAddressRequest) (*AddressResponse, error) {
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
		Address:     pgtype.Text{String: strings.TrimSpace(getStringValue(req.Address)), Valid: req.Address != nil},
		Floor:       pgtype.Text{String: strings.TrimSpace(getStringValue(req.Floor)), Valid: req.Floor != nil},
		UnitNo:      pgtype.Text{String: strings.TrimSpace(getStringValue(req.UnitNo)), Valid: req.UnitNo != nil},
		BlockTower:  pgtype.Text{String: strings.TrimSpace(getStringValue(req.BlockTower)), Valid: req.BlockTower != nil},
		CompanyName: pgtype.Text{String: strings.TrimSpace(getStringValue(req.CompanyName)), Valid: req.CompanyName != nil},
	})
	if err != nil {
		slog.Error("failed to update address for user", "user_id", userID, "address_id", addressID, "error", err)
		return nil, errors.Internal("failed to update address", err)
	}

	// Audit log the update
	s.auditService.LogUpdate(ctx, "addresses", addrID, oldAddress, address)

	return toAddressResponse(&address), nil
}

// DeleteAddress deletes an address (only if it belongs to the user)
func (s *UserService) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID string) error {
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

// SetDefaultAddress sets the default address for the authenticated user
func (s *UserService) SetDefaultAddress(ctx context.Context, userID uuid.UUID, addressID string) error {
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
		slog.Error("failed to set default address", "user_id", userID, "address_id", addressID, "error", err)
		return errors.Internal("failed to set default address", err)
	}

	return nil
}
