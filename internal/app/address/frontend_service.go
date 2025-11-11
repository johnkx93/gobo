package address

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/audit"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
)

// FrontendService handles user address operations
// Users can only manage their OWN addresses
type FrontendService struct {
	queries      *db.Queries
	auditService *audit.Service
}

func NewFrontendService(queries *db.Queries, auditService *audit.Service) *FrontendService {
	return &FrontendService{
		queries:      queries,
		auditService: auditService,
	}
}

// CreateAddress creates a new address for the authenticated user
func (s *FrontendService) CreateAddress(ctx context.Context, userID uuid.UUID, req UserCreateAddressRequest) (*AddressResponse, error) {
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

	// If the user has no default address, set this new address as default.
	userRec, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err == nil {
		if !userRec.DefaultAddressID.Valid {
			// Attempt to set default; log warning if it fails but don't fail the request
			if setErr := s.SetDefaultAddress(ctx, userID, addressID.String()); setErr != nil {
				slog.Warn("failed to set newly created address as default", "user_id", userID, "address_id", addressID, "error", setErr)
			}
		}
	} else if err != pgx.ErrNoRows {
		// Unexpected DB error while fetching user - log and continue
		slog.Warn("failed to check user's default address", "user_id", userID, "error", err)
	}

	return toAddressResponse(&address), nil
}

// GetAddress retrieves an address by ID (only if it belongs to the user)
func (s *FrontendService) GetAddress(ctx context.Context, userID uuid.UUID, addressID string) (*AddressResponse, error) {
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
func (s *FrontendService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]*AddressResponse, error) {
	addresses, err := s.queries.GetAddressesByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		slog.Error("failed to list user addresses", "user_id", userID, "error", err)
		return nil, errors.Internal("failed to list addresses", err)
	}

	// Fetch user to get DefaultAddressID for marking IsDefault
	userRec, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		// If user missing, log and continue returning addresses without defaults
		if err != pgx.ErrNoRows {
			slog.Error("failed to get user for default address check", "user_id", userID, "error", err)
			return nil, errors.Internal("failed to list addresses", err)
		}
	}

	// Sort addresses: default first, then updated_at desc
	if userRec.DefaultAddressID.Valid {
		defaultID := uuid.UUID(userRec.DefaultAddressID.Bytes)
		sort.SliceStable(addresses, func(i, j int) bool {
			ai := uuid.UUID(addresses[i].ID.Bytes)
			aj := uuid.UUID(addresses[j].ID.Bytes)

			isDefaultI := ai == defaultID
			isDefaultJ := aj == defaultID
			if isDefaultI != isDefaultJ {
				return isDefaultI // true first
			}
			// Both same default-ness: sort by updated_at desc
			return addresses[i].UpdatedAt.Time.After(addresses[j].UpdatedAt.Time)
		})
	} else {
		// No default; just sort by updated_at desc
		sort.SliceStable(addresses, func(i, j int) bool {
			return addresses[i].UpdatedAt.Time.After(addresses[j].UpdatedAt.Time)
		})
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, addr := range addresses {
		a := addr
		resp := toAddressResponse(&a)
		// Determine if this address is the user's default
		if userRec.DefaultAddressID.Valid {
			if uuid.UUID(a.ID.Bytes) == uuid.UUID(userRec.DefaultAddressID.Bytes) {
				resp.IsDefault = true
			}
		}
		responses[i] = resp
	}

	return responses, nil
}

// UpdateAddress updates an address (only if it belongs to the user)
func (s *FrontendService) UpdateAddress(ctx context.Context, userID uuid.UUID, addressID string, req UserUpdateAddressRequest) (*AddressResponse, error) {
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
func (s *FrontendService) DeleteAddress(ctx context.Context, userID uuid.UUID, addressID string) error {
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
func (s *FrontendService) SetDefaultAddress(ctx context.Context, userID uuid.UUID, addressID string) error {
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
