// Package keys provides RSA key pair management functionality.
// It implements secure generation, storage, and retrieval of RSA key pairs
// for use in cryptographic operations. The package handles key persistence
// in PEM format and provides a simple interface for key management operations.
package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrKeyGeneration = errors.New("failed to generate RSA key pair")
	ErrKeySave       = errors.New("failed to save RSA key pair")
	ErrKeyLoad       = errors.New("failed to load RSA key pair")
	ErrInvalidPath   = errors.New("invalid key path")
	ErrKeyNotFound   = errors.New("key pair not found")
)

// KeyPair represents an RSA key pair
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	KeyID      string // Unique identifier for the key pair
}

// Manager handles RSA key pair operations
type Manager struct {
	keysDir string
}

// KeysDir returns the directory where keys are stored
func (m *Manager) KeysDir() string {
	return m.keysDir
}

// NewManager creates a new key manager
func NewManager(keysDir string) (*Manager, error) {
	if keysDir == "" {
		return nil, fmt.Errorf("%w: keys directory cannot be empty", ErrInvalidPath)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, fmt.Errorf("%w: failed to create keys directory: %v", ErrKeySave, err)
	}

	return &Manager{
		keysDir: keysDir,
	}, nil
}

// GenerateKeyPair creates a new RSA key pair
func (m *Manager) GenerateKeyPair(bits int) (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		slog.Error("failed to generate RSA key pair", "error", err)
		return nil, ErrKeyGeneration
	}

	// Generate a unique key ID (using the first 8 bytes of the public key modulus)
	keyID := fmt.Sprintf("%x", privateKey.PublicKey.N.Bytes()[:8])

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		KeyID:      keyID,
	}, nil
}

// SaveKeyPair saves the RSA key pair to files
func (m *Manager) SaveKeyPair(kp *KeyPair) error {
	if kp == nil {
		return fmt.Errorf("%w: key pair cannot be nil", ErrKeySave)
	}

	privateKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.private.pem", kp.KeyID))
	publicKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.public.pem", kp.KeyID))

	// Save private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(kp.PrivateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		slog.Error("failed to save private key", "error", err)
		return ErrKeySave
	}

	// Save public key
	publicKeyBytes := x509.MarshalPKCS1PublicKey(kp.PublicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		slog.Error("failed to save public key", "error", err)
		return ErrKeySave
	}

	return nil
}

// LoadKeyPair loads an RSA key pair from files
func (m *Manager) LoadKeyPair(keyID string) (*KeyPair, error) {
	if keyID == "" {
		return nil, fmt.Errorf("%w: key ID cannot be empty", ErrKeyLoad)
	}

	privateKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.private.pem", keyID))
	publicKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.public.pem", keyID))

	// Load private key
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		slog.Error("failed to read private key", "error", err)
		return nil, ErrKeyLoad
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	if privateKeyBlock == nil {
		slog.Error("failed to decode private key PEM")
		return nil, ErrKeyLoad
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		slog.Error("failed to parse private key", "error", err)
		return nil, ErrKeyLoad
	}

	// Load public key
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		slog.Error("failed to read public key", "error", err)
		return nil, ErrKeyLoad
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		slog.Error("failed to decode public key PEM")
		return nil, ErrKeyLoad
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		slog.Error("failed to parse public key", "error", err)
		return nil, ErrKeyLoad
	}

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		KeyID:      keyID,
	}, nil
}

// ListKeyPairs returns a list of all available key pairs
func (m *Manager) ListKeyPairs() ([]string, error) {
	files, err := os.ReadDir(m.keysDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read keys directory: %w", err)
	}

	var keyIDs []string
	seen := make(map[string]bool)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Extract key ID from filename (format: {keyID}.{type}.pem)
		ext := filepath.Ext(file.Name())
		if ext != ".pem" {
			continue
		}

		base := filepath.Base(file.Name())
		parts := strings.Split(base, ".")
		if len(parts) != 3 {
			continue
		}

		keyID := parts[0]
		if !seen[keyID] {
			keyIDs = append(keyIDs, keyID)
			seen[keyID] = true
		}
	}

	return keyIDs, nil
}

// DeleteKeyPair deletes a key pair by its ID
func (m *Manager) DeleteKeyPair(keyID string) error {
	if keyID == "" {
		return fmt.Errorf("%w: key ID cannot be empty", ErrKeyLoad)
	}

	privateKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.private.pem", keyID))
	publicKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.public.pem", keyID))

	// Check if both files exist
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: private key not found", ErrKeyNotFound)
	}
	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: public key not found", ErrKeyNotFound)
	}

	// Delete private key
	if err := os.Remove(privateKeyPath); err != nil {
		slog.Error("failed to delete private key", "error", err, "path", privateKeyPath)
		return fmt.Errorf("failed to delete private key: %w", err)
	}

	// Delete public key
	if err := os.Remove(publicKeyPath); err != nil {
		slog.Error("failed to delete public key", "error", err, "path", publicKeyPath)
		return fmt.Errorf("failed to delete public key: %w", err)
	}

	return nil
}
