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

// ParsePrivateKey parses an RSA private key from its ASN.1 DER encoding.
// The keyBytes parameter should be the decoded content of a PEM block of type "RSA PRIVATE KEY".
// PEM decoding should be handled by the caller.
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
