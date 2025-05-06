package token

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// expectedTimeLag represents the maximum acceptable difference between our test's timestamp
// and the token's timestamps. Since we record the test's timestamp before the token generation,
// we expect a small lag in the token's timestamps, which is acceptable for our use case.
const expectedTimeLag = 2 // seconds

// a clear overview of all token-related test cases and their relationships. The complexity
// comes from thorough validation of JWT claims and error cases, which is essential for
// security-critical token generation.
//
//nolint:gocyclo // We intentionally keep all token generation tests in one function to maintain
func TestGenerateToken(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	generator := NewGenerator(privateKey)

	t.Run("successful token generation", func(t *testing.T) {
		username := "testuser"
		token, err := generator.GenerateToken(username)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}
		if token == "" {
			t.Error("Generated token is empty")
		}

		// Parse and verify the token
		parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		})
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}
		if !parsedToken.Valid {
			t.Error("Parsed token is invalid")
		}

		// Verify claims
		claims, maybeok := parsedToken.Claims.(jwt.MapClaims)
		if !maybeok {
			t.Fatal("Failed to parse claims")
		}

		// Check issuer
		iss, err := claims.GetIssuer()
		if err != nil {
			t.Fatal("Failed to get issuer claim")
		}
		if iss != issuerName {
			t.Errorf("Expected issuer '%s', got '%v'", issuerName, iss)
		}

		// Check username
		sub, err := claims.GetSubject()
		if err != nil {
			t.Fatal("Failed to get subject claim")
		}
		if sub != username {
			t.Errorf("Expected subject '%s', got '%v'", username, sub)
		}

		// Check timestamps
		now := time.Now().Unix()
		iat, err := claims.GetIssuedAt()
		if err != nil {
			t.Fatal("Failed to parse iat claim")
		}
		if diff := now - iat.Unix(); diff < -expectedTimeLag || diff > expectedTimeLag {
			t.Errorf("iat claim too far from now: %d seconds", diff)
		}

		nbf, err := claims.GetNotBefore()
		if err != nil {
			t.Fatal("Failed to parse nbf claim")
		}
		if diff := now - nbf.Unix(); diff < -expectedTimeLag || diff > expectedTimeLag {
			t.Errorf("nbf claim too far from now: %d seconds", diff)
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			t.Fatal("Failed to parse exp claim")
		}
		if diff := exp.Unix() - now; diff < 3600-expectedTimeLag || diff > 3600+expectedTimeLag {
			t.Errorf("exp claim not ~1 hour from now: %d seconds", diff)
		}
	})

	t.Run("empty token and error on failed signing", func(t *testing.T) {
		// Create a generator with nil private key
		invalidGenerator := NewGenerator(nil)

		token, err := invalidGenerator.GenerateToken("testuser")
		if err == nil {
			t.Error("Expected error for invalid private key")
		}
		if token != "" {
			t.Error("Expected empty token for invalid private key")
		}
	})

	t.Run("empty token and error on invalid private key", func(t *testing.T) {
		// Generate a valid key pair first
		validKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("Failed to generate key pair: %v", err)
		}

		// Create a new private key with invalid parameters
		invalidKey := &rsa.PrivateKey{
			PublicKey: validKey.PublicKey, // Keep the same public key
			D:         nil,                // Invalid private exponent
			Primes:    []*big.Int{},       // Empty primes
		}

		invalidGenerator := NewGenerator(invalidKey)
		token, err := invalidGenerator.GenerateToken("testuser")
		if err == nil {
			t.Error("Expected error for invalid private key parameters")
		}
		if token != "" {
			t.Error("Expected empty token for invalid private key parameters")
		}
	})

	// RFC 7519 Section 4.1.2 defines 'sub' (subject) as a case-sensitive string
	// that is locally unique in the context of the issuer or globally unique.
	// An empty subject would violate the uniqueness requirement.
	t.Run("empty username validation", func(t *testing.T) {
		token, err := generator.GenerateToken("")
		if err == nil {
			t.Error("Expected error for empty username")
		}
		if token != "" {
			t.Error("Expected empty token when username is empty")
		}
	})

	t.Run("token timestamps are sequential", func(t *testing.T) {
		token, err := generator.GenerateToken("testuser")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		parsedToken, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		})
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		claims, maybeok := parsedToken.Claims.(jwt.MapClaims)
		if !maybeok {
			t.Fatal("Failed to parse claims")
		}

		iat, err := claims.GetIssuedAt()
		if err != nil {
			t.Fatal("Failed to parse iat claim")
		}
		nbf, err := claims.GetNotBefore()
		if err != nil {
			t.Fatal("Failed to parse nbf claim")
		}
		exp, err := claims.GetExpirationTime()
		if err != nil {
			t.Fatal("Failed to parse exp claim")
		}

		// Verify timestamps are in correct order
		if iat.After(nbf.Time) {
			t.Error("iat should be less than or equal to nbf")
		}
		if nbf.After(exp.Time) {
			t.Error("nbf should be less than or equal to exp")
		}
		if exp.Unix()-iat.Unix() != 3600 {
			t.Errorf("Expected 1 hour between iat and exp, got %v seconds", exp.Unix()-iat.Unix())
		}
	})
}
