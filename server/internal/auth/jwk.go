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
	"oauth2-task/internal/request"
	"oauth2-task/internal/token"
)

// JWK represents a JSON Web Key as defined in RFC 7517.
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS represents a JSON Web Key Set as defined in RFC 7517.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// HandleJWKS returns the JSON Web Key Set for the server.
func HandleJWKS(keyPair token.KeyPair) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received JWKS request", "method", r.Method, "path", r.URL.Path)

		if !request.ValidateMethod(w, r, http.MethodGet) {
			return
		}

		writeJWKSResponse(w, keyPair)
	}
}

// writeJWKSResponse writes the JWKS response to the given http.ResponseWriter.
func writeJWKSResponse(w http.ResponseWriter, keyPair token.KeyPair) {
	if keyPair == nil {
		slog.Error("Invalid key pair")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_key_pair"})
		return
	}

	jwks := JWKS{
		Keys: []JWK{convertToJWK(keyPair)},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jwks); err != nil {
		slog.Error("Failed to encode JWKS response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info("Successfully sent JWKS response")
}

// convertToJWK converts a KeyPair to JWK format
func convertToJWK(keyPair token.KeyPair) JWK {
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "1", // TODO: Implement proper key ID generation
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(keyPair.PublicKey().N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(keyPair.PublicKey().E)).Bytes()),
	}
}
