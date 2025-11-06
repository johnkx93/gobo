package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Recovery middleware recovers from panics and logs the error
func Recovery(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

// Timeout middleware enforces a timeout for requests
func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.Timeout(timeout)(next)
	}
}

// ContentType middleware sets the content type to JSON
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
