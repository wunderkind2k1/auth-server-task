// Package auth implements OAuth2 authentication endpoints and JWKS functionality.
// It provides handlers for token issuance (RFC 6749), token introspection (RFC 7662),
// and JSON Web Key Set (JWKS) endpoints (RFC 7517). The package handles Basic
// Authentication, JWT token generation, token validation, and exposes public keys
// in JWKS format.
package auth

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"math/big"
	"net/http"
	"oauth2-task/internal/token"
)

// JWK represents a JSON Web Key as defined in RFC 7517
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS represents a JSON Web Key Set as defined in RFC 7517
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// HandleJWKS returns the JSON Web Key Set for the server
func HandleJWKS(keyPair token.KeyPair) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("Method not allowed", "method", r.Method, "status", http.StatusMethodNotAllowed)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Convert RSA public key to JWK format
		jwk := JWK{
			Kty: "RSA",
			Use: "sig",
			Kid: "1", // TODO: Implement proper key ID generation
			Alg: "RS256",
			N:   base64.RawURLEncoding.EncodeToString(keyPair.PublicKey().N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(keyPair.PublicKey().E)).Bytes()),
		}

		jwks := JWKS{
			Keys: []JWK{jwk},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			slog.Error("Failed to encode JWKS response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
