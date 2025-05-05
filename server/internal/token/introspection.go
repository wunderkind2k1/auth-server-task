package token

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ErrorResponse represents an OAuth2 error response.
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// IntrospectionResponse represents the OAuth2 token introspection response.
// As defined in RFC 7662 Section 2.2.
type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Username  string `json:"username,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Nbf       int64  `json:"nbf,omitempty"`
	Sub       string `json:"sub,omitempty"`
	Aud       string `json:"aud,omitempty"`
	Iss       string `json:"iss,omitempty"`
	Jti       string `json:"jti,omitempty"`
}

// validateToken parses and validates a JWT token using the provided key pair.
func validateToken(tokenString string, keyPair KeyPair) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return keyPair.PublicKey(), nil
	})
}

// extractTokenFromRequest extracts the token from either the form data or Authorization header.
func extractTokenFromRequest(r *http.Request) string {
	// Try form value first
	token := r.FormValue("token")
	if token != "" {
		return token
	}

	// Try Authorization header
	auth := r.Header.Get("Authorization")
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}

// introspectToken analyzes a validated token and returns the introspection response.
// if this gets complex, leave std lib and use. e.g. https://github.com/zitadel/zitadel
func introspectToken(parsedToken *jwt.Token) IntrospectionResponse {
	if !parsedToken.Valid {
		return IntrospectionResponse{Active: false}
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return IntrospectionResponse{Active: false}
	}

	return IntrospectionResponse{
		Active:    true,
		TokenType: "Bearer",
		Sub:       claims.Subject,
		Iss:       claims.Issuer,
		Exp:       claims.ExpiresAt.Unix(),
		Iat:       claims.IssuedAt.Unix(),
		Nbf:       claims.NotBefore.Unix(),
	}
}

func writeIntrospectionError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// For token validation failures, return active=false as per RFC 7662
	if message == "Token validation failed" {
		if err := json.NewEncoder(w).Encode(IntrospectionResponse{Active: false}); err != nil {
			slog.Error("Error encoding response", "error", err)
		}
		return
	}

	// For other errors, return the error response
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	}); err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

// HandleIntrospection processes token introspection requests.
// As defined in RFC 7662 Section 2.1.
func HandleIntrospection(keyPair KeyPair) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Technical: HTTP method validation
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			slog.Error("Method not allowed", "method", r.Method, "status", http.StatusMethodNotAllowed)
			return
		}

		// Technical: Token extraction
		tokenString := extractTokenFromRequest(r)
		if tokenString == "" {
			writeIntrospectionError(w, http.StatusBadRequest, "No token provided")
			slog.Error("No token provided for introspection")
			return
		}

		// Technical: Token validation
		parsedToken, err := validateToken(tokenString, keyPair)
		if err != nil {
			slog.Error("Token validation failed", "error", err)
			writeIntrospectionError(w, http.StatusOK, "Token validation failed")
			return
		}

		// Business Logic: Token introspection
		response := introspectToken(parsedToken)

		// Technical: Response handling
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			slog.Error("Failed to encode introspection response", "error", err)
			writeIntrospectionError(w, http.StatusInternalServerError, "Failed to encode introspection response")
			return
		}
	}
}
