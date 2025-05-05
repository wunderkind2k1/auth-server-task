// Package token provides JWT token generation, validation, and introspection functionality.
// It implements token operations using RSA private keys and handles key loading
// from files. The package follows JWT standards for token creation, signing,
// and introspection as defined in RFC 7519 and RFC 7662.
package token

import (
	"crypto/rsa"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generator handles JWT token generation.
type Generator struct {
	privateKey *rsa.PrivateKey
}

// NewGenerator creates a new token generator.
func NewGenerator(privateKey *rsa.PrivateKey) *Generator {
	return &Generator{privateKey: privateKey}
}

// GenerateToken creates a new JWT token for the given username.
func (g *Generator) GenerateToken(username string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "oauth2-server",
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		slog.Error("Failed to sign token", "error", err)
		return "", err
	}

	return tokenString, nil
}
