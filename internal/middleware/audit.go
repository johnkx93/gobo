package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/user/coc/internal/audit"
)

// RequestID middleware generates a unique request ID for each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header (from load balancer, etc.)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := audit.WithRequestID(r.Context(), requestID)

		// Add request ID to response header for tracing
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuditContext middleware extracts IP address and user agent from request
func AuditContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract IP address (handle X-Forwarded-For for proxies/load balancers)
		ipAddress := getIPAddress(r)
		ctx = audit.WithIPAddress(ctx, ipAddress)

		// Extract user agent
		userAgent := r.Header.Get("User-Agent")
		ctx = audit.WithUserAgent(ctx, userAgent)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getIPAddress extracts the real IP address from the request
// Handles X-Forwarded-For, X-Real-IP headers for reverse proxies
func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header (comma-separated list, first is original client)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Get first IP from comma-separated list
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr (format: "IP:port")
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}
