package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey reads and parses an RSA private key from a PEM file.
// Note: This function intentionally duplicates code from pkg/keys package
// as that package is meant to be a separate tooling utility for key management,
// while this is the core application logic that should be self-contained.
func LoadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	// Read the private key file
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Decode PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
