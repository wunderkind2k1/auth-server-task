package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"oauth2-task/internal/token"
	"testing"
)

// mockResponseWriter is a simple mock of http.ResponseWriter
type mockResponseWriter struct {
	headers    http.Header
	statusCode int
	body       []byte
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.body = b
	return len(b), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestHandleJWKS(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Convert private key to bytes
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPair, err := token.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	handler := HandleJWKS(keyPair)

	t.Run("rejects non-GET requests", func(t *testing.T) {
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
		for _, method := range methods {
			t.Run(method, func(t *testing.T) {
				req, err := http.NewRequest(method, "/.well-known/jwks.json", nil)
				if err != nil {
					t.Fatalf("Failed to create test request: %v", err)
				}

				w := newMockResponseWriter()
				handler(w, req)

				if w.statusCode != http.StatusMethodNotAllowed {
					t.Errorf("Expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, w.statusCode)
				}
			})
		}
	})

	t.Run("returns JWKS for GET request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
		if err != nil {
			t.Fatalf("Failed to create test request: %v", err)
		}

		w := newMockResponseWriter()
		handler(w, req)

		if w.statusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.statusCode)
		}

		if w.headers.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", w.headers.Get("Content-Type"))
		}

		var jwks JWKS
		if err := json.Unmarshal(w.body, &jwks); err != nil {
			t.Fatalf("Failed to decode JWKS response: %v", err)
		}

		if len(jwks.Keys) != 1 {
			t.Fatalf("Expected 1 key in JWKS, got %d", len(jwks.Keys))
		}

		jwk := jwks.Keys[0]
		if jwk.Kty != "RSA" {
			t.Errorf("Expected kty RSA, got %s", jwk.Kty)
		}
		if jwk.Use != "sig" {
			t.Errorf("Expected use sig, got %s", jwk.Use)
		}
		if jwk.Kid != "1" {
			t.Errorf("Expected kid 1, got %s", jwk.Kid)
		}
		if jwk.Alg != "RS256" {
			t.Errorf("Expected alg RS256, got %s", jwk.Alg)
		}
		if jwk.N == "" {
			t.Error("Expected non-empty N")
		}
		if jwk.E == "" {
			t.Error("Expected non-empty E")
		}
	})

	t.Run("handles invalid key pair", func(t *testing.T) {
		// Create a handler with nil key pair
		invalidHandler := HandleJWKS(nil)

		req, err := http.NewRequest(http.MethodGet, "/.well-known/jwks.json", nil)
		if err != nil {
			t.Fatalf("Failed to create test request: %v", err)
		}

		w := newMockResponseWriter()
		invalidHandler(w, req)

		if w.statusCode != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid key pair, got %d", http.StatusBadRequest, w.statusCode)
		}

		var errorResponse map[string]string
		if err := json.Unmarshal(w.body, &errorResponse); err != nil {
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if errorResponse["error"] != "invalid_key_pair" {
			t.Errorf("Expected error 'invalid_key_pair', got %s", errorResponse["error"])
		}
	})
}

func TestConvertToJWK(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}

	// Convert private key to bytes
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPair, err := token.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	t.Run("converts RSA public key to JWK format", func(t *testing.T) {
		jwk := convertToJWK(keyPair)

		if jwk.Kty != "RSA" {
			t.Errorf("Expected kty RSA, got %s", jwk.Kty)
		}
		if jwk.Use != "sig" {
			t.Errorf("Expected use sig, got %s", jwk.Use)
		}
		if jwk.Kid != "1" {
			t.Errorf("Expected kid 1, got %s", jwk.Kid)
		}
		if jwk.Alg != "RS256" {
			t.Errorf("Expected alg RS256, got %s", jwk.Alg)
		}
		if jwk.N == "" {
			t.Error("Expected non-empty N")
		}
		if jwk.E == "" {
			t.Error("Expected non-empty E")
		}
	})

	t.Run("returns consistent JWK for same key", func(t *testing.T) {
		jwk1 := convertToJWK(keyPair)
		jwk2 := convertToJWK(keyPair)

		if jwk1.N != jwk2.N {
			t.Error("Expected same N value for same key")
		}
		if jwk1.E != jwk2.E {
			t.Error("Expected same E value for same key")
		}
	})
}
