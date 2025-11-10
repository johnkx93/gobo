package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
)

// LogError logs an error to the error_logs table
func (s *Service) LogError(ctx context.Context, errorType, errorMessage string, stackTrace *string, requestPath, requestMethod *string) error {
	auditCtx := ExtractAuditContext(ctx)

	// If stack trace not provided, capture current stack
	var stack string
	if stackTrace != nil {
		stack = *stackTrace
	} else {
		stack = string(debug.Stack())
	}

	params := db.CreateErrorLogParams{
		UserID:        pgtype.UUID{Bytes: auditCtx.UserID, Valid: auditCtx.UserID != uuid.Nil},
		RequestID:     pgtype.Text{String: auditCtx.RequestID, Valid: auditCtx.RequestID != ""},
		ErrorType:     errorType,
		ErrorMessage:  errorMessage,
		StackTrace:    pgtype.Text{String: stack, Valid: stack != ""},
		RequestPath:   pgtype.Text{String: ptrToString(requestPath), Valid: requestPath != nil},
		RequestMethod: pgtype.Text{String: ptrToString(requestMethod), Valid: requestMethod != nil},
		IpAddress:     pgtype.Text{String: auditCtx.IPAddress, Valid: auditCtx.IPAddress != ""},
		UserAgent:     pgtype.Text{String: auditCtx.UserAgent, Valid: auditCtx.UserAgent != ""},
		Metadata:      nil,
		CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err := s.queries.CreateErrorLog(ctx, params)
	if err != nil {
		slog.Error("failed to create error log", "error", err, "error_type", errorType)
		return nil // Don't fail the main operation
	}

	return nil
}

// LogErrorWithMetadata logs an error with additional metadata
func (s *Service) LogErrorWithMetadata(ctx context.Context, errorType, errorMessage string, metadata map[string]interface{}, requestPath, requestMethod *string) error {
	auditCtx := ExtractAuditContext(ctx)

	// Convert metadata to JSON
	var metadataJSON []byte
	var err error
	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			slog.Error("failed to marshal error metadata", "error", err)
			metadataJSON = nil
		}
	}

	params := db.CreateErrorLogParams{
		UserID:        pgtype.UUID{Bytes: auditCtx.UserID, Valid: auditCtx.UserID != uuid.Nil},
		RequestID:     pgtype.Text{String: auditCtx.RequestID, Valid: auditCtx.RequestID != ""},
		ErrorType:     errorType,
		ErrorMessage:  errorMessage,
		StackTrace:    pgtype.Text{String: string(debug.Stack()), Valid: true},
		RequestPath:   pgtype.Text{String: ptrToString(requestPath), Valid: requestPath != nil},
		RequestMethod: pgtype.Text{String: ptrToString(requestMethod), Valid: requestMethod != nil},
		IpAddress:     pgtype.Text{String: auditCtx.IPAddress, Valid: auditCtx.IPAddress != ""},
		UserAgent:     pgtype.Text{String: auditCtx.UserAgent, Valid: auditCtx.UserAgent != ""},
		Metadata:      metadataJSON,
		CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err = s.queries.CreateErrorLog(ctx, params)
	if err != nil {
		slog.Error("failed to create error log", "error", err, "error_type", errorType)
		return nil // Don't fail the main operation
	}

	return nil
}

// GetRecentErrors retrieves recent errors from the error log
func (s *Service) GetRecentErrors(ctx context.Context, limit, offset int32) ([]db.ErrorLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.queries.ListRecentErrors(ctx, db.ListRecentErrorsParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		slog.Error("failed to get recent errors", "error", err)
		return nil, err
	}

	return logs, nil
}

// GetErrorsByType retrieves errors by type
func (s *Service) GetErrorsByType(ctx context.Context, errorType string, limit, offset int32) ([]db.ErrorLog, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	logs, err := s.queries.ListErrorLogsByType(ctx, db.ListErrorLogsByTypeParams{
		ErrorType: errorType,
		Limit:     int64(limit),
		Offset:    int64(offset),
	})
	if err != nil {
		slog.Error("failed to get errors by type", "error", err, "error_type", errorType)
		return nil, err
	}

	return logs, nil
}

// Helper function
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
