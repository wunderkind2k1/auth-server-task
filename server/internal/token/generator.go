// Package token provides JWT token generation, validation, and introspection functionality.
// It implements token operations using RSA private keys and handles key loading
// from files. The package follows JWT standards for token creation, signing,
// and introspection as defined in RFC 7519 and RFC 7662.
package token

import (
	"crypto/rsa"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// issuerName is the name of the token issuer as defined in RFC 7662.
const issuerName = "oauth2-server"

var (
	// ErrNilPrivateKey is returned when attempting to generate a token with a nil private key.
	ErrNilPrivateKey = errors.New("private key cannot be nil")
	// ErrEmptyUsername is returned when attempting to generate a token with an empty username.
	// RFC 7519 Section 4.1.2 requires the subject to be locally or globally unique.
	ErrEmptyUsername = errors.New("username cannot be empty")
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
	if g.privateKey == nil {
		slog.Error("Failed to validate private key", "error", ErrNilPrivateKey)
		return "", ErrNilPrivateKey
	}

	if err := g.privateKey.Validate(); err != nil {
		slog.Error("Failed to validate private key", "error", err)
		return "", err
	}

	if username == "" {
		slog.Error("Failed to generate token", "error", ErrEmptyUsername)
		return "", ErrEmptyUsername
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    issuerName,
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
