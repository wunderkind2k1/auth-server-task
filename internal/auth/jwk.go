package auth

import (
	"encoding/base64"
	"encoding/json"
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
		json.NewEncoder(w).Encode(jwks)
	}
}
