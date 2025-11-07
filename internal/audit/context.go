package audit

import (
	"context"

	"github.com/google/uuid"
)

// Context keys for audit information
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	requestIDKey contextKey = "request_id"
	ipAddressKey contextKey = "ip_address"
	userAgentKey contextKey = "user_agent"
)

// AuditContext holds audit-related information extracted from request context
type AuditContext struct {
	UserID    uuid.UUID
	RequestID string
	IPAddress string
	UserAgent string
}

// ExtractAuditContext extracts audit information from context
func ExtractAuditContext(ctx context.Context) AuditContext {
	auditCtx := AuditContext{}

	// Extract user ID
	if userID, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		auditCtx.UserID = userID
	}

	// Extract request ID
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		auditCtx.RequestID = requestID
	}

	// Extract IP address
	if ipAddress, ok := ctx.Value(ipAddressKey).(string); ok {
		auditCtx.IPAddress = ipAddress
	}

	// Extract user agent
	if userAgent, ok := ctx.Value(userAgentKey).(string); ok {
		auditCtx.UserAgent = userAgent
	}

	return auditCtx
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithIPAddress adds IP address to context
func WithIPAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, ipAddressKey, ipAddress)
}

// WithUserAgent adds user agent to context
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentKey, userAgent)
}
