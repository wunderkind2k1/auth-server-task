package token

import (
	"crypto/rsa"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generator handles JWT token generation
type Generator struct {
	privateKey *rsa.PrivateKey
}

// NewGenerator creates a new token generator
func NewGenerator(privateKey *rsa.PrivateKey) *Generator {
	return &Generator{
		privateKey: privateKey,
	}
}

// GenerateToken creates a new JWT token for the given username
func (g *Generator) GenerateToken(username string) (string, error) {
	// Create the Claims
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "auth-server",
		Subject:   username,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		slog.Error("failed to sign JWT token", "error", err)
		return "", err
	}

	return tokenString, nil
}
