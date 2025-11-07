package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
)

// Service handles audit logging
type Service struct {
	queries *db.Queries
}

// NewService creates a new audit service
func NewService(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// LogCreate logs entity creation
func (s *Service) LogCreate(ctx context.Context, entityType string, entityID uuid.UUID, newData interface{}) error {
	auditCtx := ExtractAuditContext(ctx)

	newDataJSON, err := s.prepareAuditData(newData)
	if err != nil {
		slog.Error("failed to serialize new data for audit", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	params := db.CreateAuditLogParams{
		UserID:     pgtype.UUID{Bytes: auditCtx.UserID, Valid: auditCtx.UserID != uuid.Nil},
		Action:     db.AuditActionCREATE,
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
		OldData:    nil,
		NewData:    newDataJSON,
		RequestID:  pgtype.Text{String: auditCtx.RequestID, Valid: auditCtx.RequestID != ""},
		IpAddress:  pgtype.Text{String: auditCtx.IPAddress, Valid: auditCtx.IPAddress != ""},
		UserAgent:  pgtype.Text{String: auditCtx.UserAgent, Valid: auditCtx.UserAgent != ""},
		Metadata:   nil,
		CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = s.queries.CreateAuditLog(ctx, params)
	if err != nil {
		slog.Error("failed to create audit log", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	return nil
}

// LogUpdate logs entity update
func (s *Service) LogUpdate(ctx context.Context, entityType string, entityID uuid.UUID, oldData, newData interface{}) error {
	auditCtx := ExtractAuditContext(ctx)

	oldDataJSON, err := s.prepareAuditData(oldData)
	if err != nil {
		slog.Error("failed to serialize old data for audit", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	newDataJSON, err := s.prepareAuditData(newData)
	if err != nil {
		slog.Error("failed to serialize new data for audit", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	params := db.CreateAuditLogParams{
		UserID:     pgtype.UUID{Bytes: auditCtx.UserID, Valid: auditCtx.UserID != uuid.Nil},
		Action:     db.AuditActionUPDATE,
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
		OldData:    oldDataJSON,
		NewData:    newDataJSON,
		RequestID:  pgtype.Text{String: auditCtx.RequestID, Valid: auditCtx.RequestID != ""},
		IpAddress:  pgtype.Text{String: auditCtx.IPAddress, Valid: auditCtx.IPAddress != ""},
		UserAgent:  pgtype.Text{String: auditCtx.UserAgent, Valid: auditCtx.UserAgent != ""},
		Metadata:   nil,
		CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = s.queries.CreateAuditLog(ctx, params)
	if err != nil {
		slog.Error("failed to create audit log", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	return nil
}

// LogDelete logs entity deletion
func (s *Service) LogDelete(ctx context.Context, entityType string, entityID uuid.UUID, oldData interface{}) error {
	auditCtx := ExtractAuditContext(ctx)

	oldDataJSON, err := s.prepareAuditData(oldData)
	if err != nil {
		slog.Error("failed to serialize old data for audit", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	params := db.CreateAuditLogParams{
		UserID:     pgtype.UUID{Bytes: auditCtx.UserID, Valid: auditCtx.UserID != uuid.Nil},
		Action:     db.AuditActionDELETE,
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
		OldData:    oldDataJSON,
		NewData:    nil,
		RequestID:  pgtype.Text{String: auditCtx.RequestID, Valid: auditCtx.RequestID != ""},
		IpAddress:  pgtype.Text{String: auditCtx.IPAddress, Valid: auditCtx.IPAddress != ""},
		UserAgent:  pgtype.Text{String: auditCtx.UserAgent, Valid: auditCtx.UserAgent != ""},
		Metadata:   nil,
		CreatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = s.queries.CreateAuditLog(ctx, params)
	if err != nil {
		slog.Error("failed to create audit log", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil // Don't fail the main operation
	}

	return nil
}

// GetEntityHistory retrieves audit history for a specific entity
func (s *Service) GetEntityHistory(ctx context.Context, entityType string, entityID uuid.UUID, limit, offset int32) ([]db.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.queries.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: entityType,
		EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		slog.Error("failed to get entity audit history", "error", err, "entity_type", entityType, "entity_id", entityID)
		return nil, err
	}

	return logs, nil
}

// GetUserAuditHistory retrieves audit history for a specific user
func (s *Service) GetUserAuditHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]db.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.queries.ListAuditLogsByUser(ctx, db.ListAuditLogsByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		slog.Error("failed to get user audit history", "error", err, "user_id", userID)
		return nil, err
	}

	return logs, nil
}

// prepareAuditData converts data to JSON and filters sensitive information
func (s *Service) prepareAuditData(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Parse to map to filter sensitive fields
	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return nil, err
	}

	// Filter sensitive fields
	dataMap = s.filterSensitiveData(dataMap)

	// Convert back to JSON
	return json.Marshal(dataMap)
}

// filterSensitiveData removes sensitive fields from audit data
func (s *Service) filterSensitiveData(data map[string]interface{}) map[string]interface{} {
	sensitiveFields := []string{
		"password",
		"password_hash",
		"passwordHash",
		"token",
		"secret",
		"api_key",
		"apiKey",
		"private_key",
		"privateKey",
	}

	for _, field := range sensitiveFields {
		delete(data, field)
	}

	return data
}
