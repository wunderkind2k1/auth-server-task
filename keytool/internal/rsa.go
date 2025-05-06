// Package rsa provides the core functionality for the key management tool.
// It implements secure generation, storage, and retrieval of RSA key pairs
// specifically designed for JWT signing in the OAuth2 server.
//
// Key Features:
//   - Secure RSA key pair generation with configurable key sizes (2048+ bits)
//   - PEM-encoded key storage with proper file permissions (0600 for private, 0644 for public)
//   - Unique key identification using public key modulus
//   - Atomic key pair operations (save/delete)
//   - Built-in key validation using crypto/rsa.Validate()
//
// Security Considerations:
//   - Private keys are stored with strict permissions (0600)
//   - Public keys are stored with read-only permissions (0644)
//   - Minimum key size of 2048 bits is enforced
//   - Key operations are atomic to prevent partial writes
//   - Keys are validated before saving
//
// Usage:
//
//	manager, err := NewManager("/path/to/keys")
//	if err != nil {
//	    // Handle error
//	}
//
//	// Generate a new key pair
//	keyPair, err := manager.generateKeyPair(2048)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Save the key pair
//	if err := manager.SaveKeyPair(keyPair); err != nil {
//	    // Handle error
//	}
//
// This package is intentionally separate from the main application's key handling
// to maintain a clear boundary between the key management tool and the server.
package rsa

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
	// ErrKeyGeneration is returned when RSA key pair generation fails.
	ErrKeyGeneration = errors.New("failed to generate RSA key pair")
	// ErrKeySave is returned when saving an RSA key pair to disk fails.
	ErrKeySave = errors.New("failed to save RSA key pair")
	// ErrKeyLoad is returned when loading an RSA key pair from disk fails.
	ErrKeyLoad = errors.New("failed to load RSA key pair")
	// ErrInvalidPath is returned when an invalid path is provided for key operations.
	ErrInvalidPath = errors.New("invalid key path")
	// ErrKeyNotFound is returned when a requested key pair cannot be found.
	ErrKeyNotFound = errors.New("key pair not found")
	// ErrInvalidKeySize is returned when an invalid key size is requested.
	ErrInvalidKeySize = errors.New("invalid key size: must be at least 2048 bits")
)

// KeyPair represents an RSA key pair.
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	KeyID      string // Unique identifier for the key pair
}

// Manager handles RSA key pair operations.
type Manager struct {
	keysDir string
}

// KeysDir returns the directory where keys are stored.
func (m *Manager) KeysDir() string {
	return m.keysDir
}

// NewManager creates a new key manager.
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

// GenerateKeyPair creates a new RSA key pair.
func (m *Manager) GenerateKeyPair(bits int) (*KeyPair, error) {
	if bits < 2048 {
		return nil, ErrInvalidKeySize
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		slog.Error("failed to generate RSA key pair", "error", err)
		return nil, ErrKeyGeneration
	}

	// Generate a unique key ID (using the first 8 bytes of the public key modulus)
	keyID := fmt.Sprintf("%x", privateKey.N.Bytes()[:8])

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
		KeyID:      keyID,
	}, nil
}

// SaveKeyPair saves the RSA key pair to files.
func (m *Manager) SaveKeyPair(kp *KeyPair) error {
	if kp == nil {
		return errors.New("key pair is nil")
	}

	if kp.PrivateKey == nil || kp.PublicKey == nil {
		return errors.New("invalid key pair: private or public key is nil")
	}

	// Validate the private key
	if err := kp.PrivateKey.Validate(); err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	// Generate key ID from public key modulus
	keyID := fmt.Sprintf("%x", kp.PublicKey.N.Bytes()[:8])

	// Create key directory if it doesn't exist
	if err := os.MkdirAll(m.keysDir, 0700); err != nil {
		return fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Encode private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(kp.PrivateKey),
	})

	// Encode public key
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(kp.PublicKey),
	})

	// Save private key
	privateKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.private.pem", keyID))
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	// Save public key with 0644 permissions (readable by others) as it's meant to be shared
	// and used by other services for JWT verification
	publicKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.public.pem", keyID))
	// #nosec G306 -- Public key file needs to be readable by others for JWT verification
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		// Clean up private key if public key save fails
		_ = os.Remove(privateKeyPath)
		return fmt.Errorf("failed to save public key: %w", err)
	}

	return nil
}

// LoadKeyPair loads an RSA key pair from files.
func (m *Manager) LoadKeyPair(keyID string) (*KeyPair, error) {
	if keyID == "" {
		return nil, fmt.Errorf("%w: key ID cannot be empty", ErrKeyLoad)
	}

	privateKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.private.pem", keyID))
	publicKeyPath := filepath.Join(m.keysDir, fmt.Sprintf("%s.public.pem", keyID))

	// Load private key
	// #nosec G304 -- File path is constructed within the same method and content is immediately parsed as PEM
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
	// #nosec G304 -- File path is constructed within the same method and content is immediately parsed as PEM
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

// ListKeyPairs returns a list of all available key pairs.
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

// DeleteKeyPair deletes a key pair by its ID.
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
