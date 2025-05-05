package token

import (
	"crypto/rsa"
	"crypto/x509"
	"log/slog"
)

// KeyPair defines the interface for a pair of cryptographic keys.
type KeyPair interface {
	// PrivateKey returns the private key.
	PrivateKey() *rsa.PrivateKey
	// PublicKey returns the public key.
	PublicKey() *rsa.PublicKey
}

// rsaKeyPair implements KeyPair for RSA keys.
type rsaKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// PrivateKey returns the RSA private key.
func (k *rsaKeyPair) PrivateKey() *rsa.PrivateKey {
	return k.privateKey
}

// PublicKey returns the RSA public key.
func (k *rsaKeyPair) PublicKey() *rsa.PublicKey {
	return k.publicKey
}

// ParsePrivateKey parses a PEM-encoded RSA private key.
func ParsePrivateKey(keyBytes []byte) (KeyPair, error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
	if err != nil {
		slog.Error("Failed to parse private key", "error", err)
		return nil, err
	}
	return &rsaKeyPair{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}
