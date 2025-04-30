package auth

import (
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrMissingHeader     = errors.New("authorization header required")
	ErrInvalidFormat     = errors.New("invalid authorization header format")
	ErrInvalidBase64     = errors.New("invalid base64 encoding in credentials")
	ErrInvalidAuthScheme = errors.New("invalid authorization scheme")
)

// BasicAuth represents basic authentication credentials
type BasicAuth struct {
	Username string
	Password string
}

// ErrorResponse represents an authentication error response
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GetErrorResponse returns the appropriate error response for an auth error
func GetErrorResponse(err error) ErrorResponse {
	switch err {
	case ErrMissingHeader:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Authorization header required",
		}
	case ErrInvalidFormat:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Invalid authorization header format",
		}
	case ErrInvalidBase64:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Invalid credentials",
		}
	default:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Authentication failed",
		}
	}
}

// ParseBasicAuth validates the Authorization header for Basic Auth
func ParseBasicAuth(authHeader string) error {
	if authHeader == "" {
		return ErrMissingHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return ErrInvalidFormat
	}

	// Decode the base64-encoded credentials
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return ErrInvalidBase64
	}

	// Validate credentials format (username:password)
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return ErrInvalidFormat
	}

	return nil
}
