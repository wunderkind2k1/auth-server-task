package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// KeyPair defines the interface for a pair of cryptographic keys
type KeyPair interface {
	// PrivateKey returns the private key
	PrivateKey() *rsa.PrivateKey
	// PublicKey returns the public key
	PublicKey() *rsa.PublicKey
}

// rsaKeyPair implements KeyPair for RSA keys
type rsaKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// PrivateKey returns the RSA private key
func (k *rsaKeyPair) PrivateKey() *rsa.PrivateKey {
	return k.privateKey
}

// PublicKey returns the RSA public key
func (k *rsaKeyPair) PublicKey() *rsa.PublicKey {
	return k.publicKey
}

// LoadPrivateKey reads and parses an RSA private key from a PEM file.
// Note: This function intentionally duplicates code from pkg/keys package
// as that package is meant to be a separate tooling utility for key management,
// while this is the core application logic that should be self-contained.
func LoadPrivateKey(filePath string) (KeyPair, error) {
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

	return &rsaKeyPair{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}
