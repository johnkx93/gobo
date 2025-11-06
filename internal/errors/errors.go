package errors

import "fmt"

// DomainError represents application-specific errors
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	CodeNotFound      = "NOT_FOUND"
	CodeAlreadyExists = "ALREADY_EXISTS"
	CodeValidation    = "VALIDATION_ERROR"
	CodeInternal      = "INTERNAL_ERROR"
	CodeUnauthorized  = "UNAUTHORIZED"
)

// Common error constructors
func NotFound(message string) *DomainError {
	return &DomainError{
		Code:    CodeNotFound,
		Message: message,
	}
}

func AlreadyExists(message string) *DomainError {
	return &DomainError{
		Code:    CodeAlreadyExists,
		Message: message,
	}
}

func Validation(message string) *DomainError {
	return &DomainError{
		Code:    CodeValidation,
		Message: message,
	}
}

func Internal(message string, err error) *DomainError {
	return &DomainError{
		Code:    CodeInternal,
		Message: message,
		Err:     err,
	}
}

func Unauthorized(message string) *DomainError {
	return &DomainError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}
