package token

import (
	"crypto/rsa"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims represents the JWT claims
type CustomClaims struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

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

// GenerateToken creates a new JWT token
func (g *Generator) GenerateToken() (string, error) {
	// Create the Claims
	claims := CustomClaims{
		"1234567890",
		"John Doe",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-server",
			Subject:   "1234567890",
		},
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
