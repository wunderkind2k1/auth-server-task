package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"oauth2-task/internal/token"
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

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// tokenError represents an OAuth2 error response
type tokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// HandleToken processes OAuth2 token requests
func HandleToken(keyPair token.KeyPair, userPool map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			slog.Error("Method not allowed", "method", r.Method, "status", http.StatusMethodNotAllowed)
			return
		}

		// Create BasicAuth instance with user pool
		basicAuth := NewBasicAuth(userPool)

		// Validate Basic Auth
		if err := basicAuth.ParseBasicAuth(r.Header.Get("Authorization")); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := GetErrorResponse(err)
			json.NewEncoder(w).Encode(errorResponse)
			slog.Error("Authentication failed", "error", err)
			return
		}

		// Create token generator
		generator := token.NewGenerator(keyPair.PrivateKey())

		// Generate a real JWT token
		tokenString, err := generator.GenerateToken(basicAuth.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{
				Error:            "server_error",
				ErrorDescription: "Failed to generate token",
			})
			if err != nil {
				slog.Error("Failed to generate Error Response for token generation", "error", err)
				return
			}
			slog.Error("Failed to generate token", "error", err)
			return
		}

		// Return the token response
		response := TokenResponse{
			AccessToken: tokenString,
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
