package auth

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"strings"
)

var (
	ErrMissingHeader      = errors.New("authorization header required")
	ErrInvalidFormat      = errors.New("invalid authorization header format")
	ErrInvalidBase64      = errors.New("invalid base64 encoding in credentials")
	ErrInvalidAuthScheme  = errors.New("invalid authorization scheme")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// userPool represents a collection of users and their credentials
type userPool map[string]string

// BasicAuth represents basic authentication credentials
type BasicAuth struct {
	Username string
	Password string
	pool     map[string]string
}

// NewBasicAuth creates a new BasicAuth instance with a user pool
func NewBasicAuth(userpool map[string]string) *BasicAuth {
	return &BasicAuth{
		pool: userpool,
	}
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
	case ErrInvalidCredentials:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Invalid username or password",
		}
	default:
		return ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Authentication failed",
		}
	}
}

// ParseBasicAuth validates the Authorization header for Basic Auth
func (ba *BasicAuth) ParseBasicAuth(authHeader string) error {
	if authHeader == "" {
		slog.Error(ErrMissingHeader.Error())
		return ErrMissingHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		slog.Error(ErrInvalidFormat.Error(), "header", authHeader)
		return ErrInvalidFormat
	}

	// Decode the base64-encoded credentials
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		slog.Error(ErrInvalidBase64.Error(), "error", err)
		return ErrInvalidBase64
	}

	// Validate credentials format (username:password)
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		slog.Error(ErrInvalidFormat.Error())
		return ErrInvalidFormat
	}

	// Validate credentials against the user pool
	storedPassword, exists := ba.pool[credentials[0]]
	if !exists || storedPassword != credentials[1] {
		slog.Error(ErrInvalidCredentials.Error(), "username", credentials[0])
		return ErrInvalidCredentials
	}

	// Store the validated credentials
	ba.Username = credentials[0]
	ba.Password = credentials[1]

	return nil
}
