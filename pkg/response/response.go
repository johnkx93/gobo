package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/user/coc/pkg/errors"
)

// JSONResponse represents a standard JSON response
type JSONResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON writes a JSON response with the given status code, message, and data
func JSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  status >= 200 && status < 300,
		Message: message,
		Data:    data,
	})
}

// Error writes a JSON error response with the given status code and message
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{
		Status:  false,
		Message: message,
	})
}

// HandleServiceError handles domain errors and writes appropriate HTTP responses
func HandleServiceError(w http.ResponseWriter, err error) {
	var domainErr *errors.DomainError
	if e, ok := err.(*errors.DomainError); ok {
		domainErr = e
	} else {
		slog.Error("unexpected error type", "error", err)
		Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	switch domainErr.Code {
	case errors.CodeNotFound:
		Error(w, http.StatusNotFound, domainErr.Message)
	case errors.CodeAlreadyExists:
		Error(w, http.StatusConflict, domainErr.Message)
	case errors.CodeValidation:
		Error(w, http.StatusBadRequest, domainErr.Message)
	case errors.CodeUnauthorized:
		Error(w, http.StatusUnauthorized, domainErr.Message)
	default:
		slog.Error("internal error", "error", domainErr)
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
