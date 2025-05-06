package rsa

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temporary directory for testing and returns its path
// and a cleanup function.
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "keytool-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir, func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("Failed to cleanup temp dir: %v", err)
		}
	}
}

// createTestManager creates a new Manager instance in a temporary directory.
func createTestManager(t *testing.T) (*Manager, func()) {
	t.Helper()
	dir, cleanup := setupTestDir(t)
	manager, err := NewManager(dir)
	if err != nil {
		cleanup()
		t.Fatalf("Failed to create manager: %v", err)
	}
	return manager, cleanup
}

// validateKeyPair performs basic sanity checks on the RSA key pair using the built-in Validate() method.
func validateKeyPair(t *testing.T, keyPair *KeyPair, wantBits int) {
	t.Helper()

	// Validate private key
	if err := keyPair.PrivateKey.Validate(); err != nil {
		t.Errorf("Private key validation failed: %v", err)
	}

	// Verify key sizes match requested bits
	if keyPair.PrivateKey.N.BitLen() != wantBits {
		t.Errorf("Private key size = %d bits, want %d bits", keyPair.PrivateKey.N.BitLen(), wantBits)
	}
	if keyPair.PublicKey.N.BitLen() != wantBits {
		t.Errorf("Public key size = %d bits, want %d bits", keyPair.PublicKey.N.BitLen(), wantBits)
	}
}

func TestManager_GenerateKeyPair(t *testing.T) {
	tests := []struct {
		name    string
		bits    int
		wantErr bool
	}{
		{"valid 2048 bits", 2048, false},
		{"valid 4096 bits", 4096, false},
		{"invalid bits", 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, cleanup := createTestManager(t)
			defer cleanup()

			keyPair, err := manager.GenerateKeyPair(tt.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("generateKeyPair() returned nil keyPair")
					return
				}
				if keyPair.KeyID == "" {
					t.Error("generateKeyPair() returned empty KeyID")
				}
				if keyPair.PrivateKey == nil {
					t.Error("generateKeyPair() returned nil PrivateKey")
				}
				if keyPair.PublicKey == nil {
					t.Error("generateKeyPair() returned nil PublicKey")
				}

				// Validate the generated key pair
				validateKeyPair(t, keyPair, tt.bits)
			}
		})
	}
}

func TestManager_SaveKeyPair(t *testing.T) {
	manager, cleanup := createTestManager(t)
	defer cleanup()

	// Generate a valid key pair for testing
	keyPair, err := manager.GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	tests := []struct {
		name    string
		keyPair *KeyPair
		wantErr bool
	}{
		{"valid key pair", keyPair, false},
		{"nil key pair", nil, true},
		{"invalid key pair", &KeyPair{KeyID: "test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SaveKeyPair(tt.keyPair)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify files were created with correct permissions
				privateKeyPath := filepath.Join(manager.KeysDir(), tt.keyPair.KeyID+".private.pem")
				publicKeyPath := filepath.Join(manager.KeysDir(), tt.keyPair.KeyID+".public.pem")

				// Check private key file
				if info, err := os.Stat(privateKeyPath); err != nil {
					t.Errorf("Private key file not found: %v", err)
				} else if info.Mode() != 0600 {
					t.Errorf("Private key file has wrong permissions: %v", info.Mode())
				}

				// Check public key file
				if info, err := os.Stat(publicKeyPath); err != nil {
					t.Errorf("Public key file not found: %v", err)
				} else if info.Mode() != 0644 {
					t.Errorf("Public key file has wrong permissions: %v", info.Mode())
				}
			}
		})
	}
}

func TestManager_ListKeyPairs_WithKeys(t *testing.T) {
	manager, cleanup := createTestManager(t)
	defer cleanup()

	// Generate and save two key pairs
	keyPair1, err := manager.GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate first test key pair: %v", err)
	}
	if err := manager.SaveKeyPair(keyPair1); err != nil {
		t.Fatalf("Failed to save first test key pair: %v", err)
	}

	keyPair2, err := manager.GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate second test key pair: %v", err)
	}
	if err := manager.SaveKeyPair(keyPair2); err != nil {
		t.Fatalf("Failed to save second test key pair: %v", err)
	}

	// List key pairs
	gotIDs, err := manager.ListKeyPairs()
	if err != nil {
		t.Errorf("ListKeyPairs() error = %v", err)
		return
	}

	// Verify we got both key IDs
	if len(gotIDs) != 2 {
		t.Errorf("ListKeyPairs() returned %d IDs, want 2", len(gotIDs))
		return
	}

	// Create a map for easier comparison
	gotMap := make(map[string]bool)
	for _, id := range gotIDs {
		gotMap[id] = true
	}

	if !gotMap[keyPair1.KeyID] {
		t.Errorf("ListKeyPairs() missing first key ID: %s", keyPair1.KeyID)
	}
	if !gotMap[keyPair2.KeyID] {
		t.Errorf("ListKeyPairs() missing second key ID: %s", keyPair2.KeyID)
	}
}

func TestManager_ListKeyPairs_Empty(t *testing.T) {
	manager, cleanup := createTestManager(t)
	defer cleanup()

	gotIDs, err := manager.ListKeyPairs()
	if err != nil {
		t.Errorf("ListKeyPairs() error = %v", err)
		return
	}

	if len(gotIDs) != 0 {
		t.Errorf("ListKeyPairs() returned %d IDs, want 0", len(gotIDs))
	}
}

func TestManager_DeleteKeyPair(t *testing.T) {
	manager, cleanup := createTestManager(t)
	defer cleanup()

	// Generate and save a test key pair
	keyPair, err := manager.GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}
	if err := manager.SaveKeyPair(keyPair); err != nil {
		t.Fatalf("Failed to save test key pair: %v", err)
	}

	tests := []struct {
		name    string
		keyID   string
		wantErr bool
	}{
		{"existing key", keyPair.KeyID, false},
		{"non-existent key", "nonexistent", true},
		{"empty key ID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.DeleteKeyPair(tt.keyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify files were deleted
				privateKeyPath := filepath.Join(manager.KeysDir(), tt.keyID+".private.pem")
				publicKeyPath := filepath.Join(manager.KeysDir(), tt.keyID+".public.pem")

				if _, err := os.Stat(privateKeyPath); !os.IsNotExist(err) {
					t.Error("Private key file still exists after deletion")
				}
				if _, err := os.Stat(publicKeyPath); !os.IsNotExist(err) {
					t.Error("Public key file still exists after deletion")
				}
			}
		})
	}
}
