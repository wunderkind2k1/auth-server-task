package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
)

func TestParsePrivateKey(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Get the DER encoding
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	t.Run("successful key parsing", func(t *testing.T) {
		keyPair, err := ParsePrivateKey(privateKeyBytes)
		if err != nil {
			t.Fatalf("Failed to parse valid private key: %v", err)
		}

		// Verify the interface implementation
		if keyPair == nil {
			t.Fatal("KeyPair is nil")
		}

		// Verify the concrete type
		rsaPair, ok := keyPair.(*rsaKeyPair)
		if !ok {
			t.Fatal("Expected *rsaKeyPair implementation")
		}
		if rsaPair.privateKey == nil {
			t.Fatal("rsaKeyPair.privateKey is nil")
		}
		if rsaPair.publicKey == nil {
			t.Fatal("rsaKeyPair.publicKey is nil")
		}

		// Verify private key
		parsedPrivateKey := keyPair.PrivateKey()
		if parsedPrivateKey == nil {
			t.Fatal("Private key is nil")
		}
		if parsedPrivateKey.N.Cmp(privateKey.N) != 0 {
			t.Error("Private key modulus mismatch")
		}
		if parsedPrivateKey.E != privateKey.E {
			t.Error("Private key public exponent mismatch")
		}

		// Verify public key
		parsedPublicKey := keyPair.PublicKey()
		if parsedPublicKey == nil {
			t.Fatal("Public key is nil")
		}
		if parsedPublicKey.N.Cmp(privateKey.PublicKey.N) != 0 {
			t.Error("Public key modulus mismatch")
		}
		if parsedPublicKey.E != privateKey.PublicKey.E {
			t.Error("Public key exponent mismatch")
		}
	})

	t.Run("invalid key format", func(t *testing.T) {
		invalidKeyBytes := []byte("not a valid RSA private key")
		keyPair, err := ParsePrivateKey(invalidKeyBytes)
		if err == nil {
			t.Error("Expected error for invalid key format")
		}
		if keyPair != nil {
			t.Error("Expected nil KeyPair for invalid key format")
		}
	})

	t.Run("empty key data", func(t *testing.T) {
		keyPair, err := ParsePrivateKey(nil)
		if err == nil {
			t.Error("Expected error for nil key data")
		}
		if keyPair != nil {
			t.Error("Expected nil KeyPair for nil key data")
		}
	})
}
