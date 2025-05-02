package auth

import (
	"encoding/json"
	"net/http"
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
func HandleJWKS(w http.ResponseWriter, r *http.Request) {
	// For now, return an empty key set
	// TODO: Implement actual key retrieval from the key store
	jwks := JWKS{
		Keys: []JWK{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
