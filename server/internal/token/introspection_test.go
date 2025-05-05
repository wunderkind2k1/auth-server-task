package token

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// setupTestKeyPair creates a test RSA key pair for testing
func setupTestKeyPair(t *testing.T) KeyPair {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test key pair: %v", err)
	}
	return &rsaKeyPair{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// createTestToken creates a JWT token with the given claims
func createTestToken(t *testing.T, keyPair KeyPair, claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(keyPair.PrivateKey())
	if err != nil {
		t.Fatalf("Failed to sign test token: %v", err)
	}
	return tokenString
}

func TestIntrospectToken(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		token   *jwt.Token
		want    IntrospectionResponse
		wantErr bool
	}{
		{
			name: "Valid token",
			token: &jwt.Token{
				Claims: &jwt.RegisteredClaims{
					Issuer:    "test-issuer",
					Subject:   "test-subject",
					IssuedAt:  jwt.NewNumericDate(now),
					NotBefore: jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				},
				Valid: true,
			},
			want: IntrospectionResponse{
				Active:    true,
				TokenType: "Bearer",
				Sub:       "test-subject",
				Iss:       "test-issuer",
				Exp:       now.Add(time.Hour).Unix(),
				Iat:       now.Unix(),
				Nbf:       now.Unix(),
			},
			wantErr: false,
		},
		{
			name: "Invalid token - not valid",
			token: &jwt.Token{
				Claims: &jwt.RegisteredClaims{},
				Valid:  false,
			},
			want:    IntrospectionResponse{Active: false},
			wantErr: false,
		},
		{
			name: "Invalid token - wrong claims type",
			token: &jwt.Token{
				Claims: jwt.MapClaims{}, // Different claims type
				Valid:  true,
			},
			want:    IntrospectionResponse{Active: false},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := introspectToken(tt.token)

			// Check active status
			if got.Active != tt.want.Active {
				t.Errorf("introspectToken() active = %v, want %v", got.Active, tt.want.Active)
			}

			// For active tokens, check other fields
			if got.Active {
				if got.TokenType != tt.want.TokenType {
					t.Errorf("introspectToken() tokenType = %v, want %v", got.TokenType, tt.want.TokenType)
				}
				if got.Sub != tt.want.Sub {
					t.Errorf("introspectToken() sub = %v, want %v", got.Sub, tt.want.Sub)
				}
				if got.Iss != tt.want.Iss {
					t.Errorf("introspectToken() iss = %v, want %v", got.Iss, tt.want.Iss)
				}
				if got.Exp != tt.want.Exp {
					t.Errorf("introspectToken() exp = %v, want %v", got.Exp, tt.want.Exp)
				}
				if got.Iat != tt.want.Iat {
					t.Errorf("introspectToken() iat = %v, want %v", got.Iat, tt.want.Iat)
				}
				if got.Nbf != tt.want.Nbf {
					t.Errorf("introspectToken() nbf = %v, want %v", got.Nbf, tt.want.Nbf)
				}
			}
		})
	}
}

func TestExtractTokenFromRequest(t *testing.T) {
	tests := []struct {
		name        string
		formToken   string
		authHeader  string
		wantToken   string
		description string
	}{
		{
			name:        "Token in form data",
			formToken:   "form.token.here",
			authHeader:  "",
			wantToken:   "form.token.here",
			description: "Should extract token from form data",
		},
		{
			name:        "Token in Authorization header",
			formToken:   "",
			authHeader:  "Bearer header.token.here",
			wantToken:   "header.token.here",
			description: "Should extract token from Authorization header",
		},
		{
			name:        "Form data takes precedence",
			formToken:   "form.token.here",
			authHeader:  "Bearer header.token.here",
			wantToken:   "form.token.here",
			description: "Form data should take precedence over Authorization header",
		},
		{
			name:        "No token provided",
			formToken:   "",
			authHeader:  "",
			wantToken:   "",
			description: "Should return empty string when no token is provided",
		},
		{
			name:        "Invalid Authorization header - no Bearer prefix",
			formToken:   "",
			authHeader:  "header.token.here",
			wantToken:   "",
			description: "Should return empty string when Authorization header has no Bearer prefix",
		},
		{
			name:        "Invalid Authorization header - wrong prefix",
			formToken:   "",
			authHeader:  "Basic header.token.here",
			wantToken:   "",
			description: "Should return empty string when Authorization header has wrong prefix",
		},
		{
			name:        "Invalid Authorization header - empty after Bearer",
			formToken:   "",
			authHeader:  "Bearer ",
			wantToken:   "",
			description: "Should return empty string when Authorization header has no token after Bearer",
		},
		{
			name:        "Invalid Authorization header - Bearer with spaces",
			formToken:   "",
			authHeader:  "Bearer  header.token.here",
			wantToken:   " header.token.here",
			description: "Should preserve spaces after Bearer prefix",
		},
		{
			name:        "Invalid Authorization header - Bearer with tabs",
			formToken:   "",
			authHeader:  "Bearer\theader.token.here",
			wantToken:   "",
			description: "Should return empty string when Bearer is followed by tab",
		},
		{
			name:        "Invalid Authorization header - Bearer with newline",
			formToken:   "",
			authHeader:  "Bearer\nheader.token.here",
			wantToken:   "",
			description: "Should return empty string when Bearer is followed by newline",
		},
		{
			name:        "Invalid Authorization header - Bearer with carriage return",
			formToken:   "",
			authHeader:  "Bearer\rheader.token.here",
			wantToken:   "",
			description: "Should return empty string when Bearer is followed by carriage return",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			req, err := http.NewRequest("POST", "/introspect", nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %v", err)
			}

			// Set form data if provided
			if tt.formToken != "" {
				req.Form = url.Values{}
				req.Form.Set("token", tt.formToken)
			}

			// Set Authorization header if provided
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Test token extraction
			got := extractTokenFromRequest(req)
			if got != tt.wantToken {
				t.Errorf("extractTokenFromRequest() = %v, want %v (%s)", got, tt.wantToken, tt.description)
			}
		})
	}
}

func TestValidateSigningMethod(t *testing.T) {
	// Setup test key pair
	keyPair := setupTestKeyPair(t)

	tests := []struct {
		name      string
		token     *jwt.Token
		wantKey   interface{}
		wantError error
	}{
		{
			name: "Valid RSA token",
			token: &jwt.Token{
				Method: jwt.SigningMethodRS256,
			},
			wantKey:   keyPair.PublicKey(),
			wantError: nil,
		},
		{
			name: "Invalid signing method - HMAC",
			token: &jwt.Token{
				Method: jwt.SigningMethodHS256,
			},
			wantKey:   nil,
			wantError: jwt.ErrSignatureInvalid,
		},
		{
			name: "Invalid signing method - ECDSA",
			token: &jwt.Token{
				Method: jwt.SigningMethodES256,
			},
			wantKey:   nil,
			wantError: jwt.ErrSignatureInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotErr := validateSigningMethod(tt.token, keyPair)

			// Check error
			if gotErr != tt.wantError {
				t.Errorf("validateSigningMethod() error = %v, want %v", gotErr, tt.wantError)
			}

			// Check key
			if tt.wantError == nil {
				if gotKey != tt.wantKey {
					t.Errorf("validateSigningMethod() key = %v, want %v", gotKey, tt.wantKey)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	keyPair := setupTestKeyPair(t)
	now := time.Now()

	// Create a valid token
	validToken := createTestToken(t, keyPair, jwt.RegisteredClaims{
		Issuer:    "test-issuer",
		Subject:   "test-subject",
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	})

	tests := []struct {
		name      string
		token     string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "Valid token",
			token:     validToken,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "Invalid token - malformed",
			token:     "invalid.token.string",
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Invalid token - empty",
			token:     "",
			wantValid: false,
			wantErr:   true,
		},
		{
			name:      "Invalid token - wrong signature",
			token:     validToken + "tampered",
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateToken(tt.token, keyPair)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check token validity
			if got != nil && got.Valid != tt.wantValid {
				t.Errorf("validateToken() valid = %v, want %v", got.Valid, tt.wantValid)
			}
		})
	}
}

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

// TestWriteIntrospectionError tests the error response writing functionality.
// While it's a general test for writeIntrospectionError, it has two critical business requirements:
//
//  1. Status Code Preservation:
//     The status code passed to writeIntrospectionError must be preserved in the response.
//     This is important because different error scenarios require different HTTP status codes,
//     and we must ensure they are correctly propagated to the client.
//
//  2. Token Validation Failure Response:
//     When a token validation fails (message == "Token validation failed"),
//     the response must always be IntrospectionResponse{Active: false}.
//     This is a critical security requirement as per RFC 7662, ensuring that
//     invalid tokens are properly marked as inactive without exposing internal details.
func TestWriteIntrospectionError(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		status        int
		wantResponse  IntrospectionResponse
		wantErrorResp ErrorResponse
		wantStatus    int
	}{
		{
			name:    "Token validation failed",
			message: "Token validation failed",
			status:  http.StatusOK,
			wantResponse: IntrospectionResponse{
				Active: false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "Other error",
			message: "Some other error",
			status:  http.StatusBadRequest,
			wantErrorResp: ErrorResponse{
				Error: "Some other error",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newMockResponseWriter()

			// Call the function
			writeIntrospectionError(w, tt.status, tt.message)

			// Check status code
			if w.statusCode != tt.wantStatus {
				t.Errorf("writeIntrospectionError() status = %v, want %v", w.statusCode, tt.wantStatus)
			}

			// Check content type
			if w.headers.Get("Content-Type") != "application/json" {
				t.Errorf("writeIntrospectionError() content-type = %v, want application/json", w.headers.Get("Content-Type"))
			}

			// Check response body
			if tt.message == "Token validation failed" {
				var got IntrospectionResponse
				if err := json.Unmarshal(w.body, &got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if got != tt.wantResponse {
					t.Errorf("writeIntrospectionError() response = %v, want %v", got, tt.wantResponse)
				}
			} else {
				var got ErrorResponse
				if err := json.Unmarshal(w.body, &got); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if got != tt.wantErrorResp {
					t.Errorf("writeIntrospectionError() error response = %v, want %v", got, tt.wantErrorResp)
				}
			}
		})
	}
}

// TestValidateHTTPMethod tests the HTTP method validation as required by RFC 7662.
// The introspection endpoint must only accept POST requests.
func TestValidateHTTPMethod(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid POST method",
			method:     http.MethodPost,
			wantStatus: 0, // No status written for valid method
			wantBody:   "",
		},
		{
			name:       "Invalid GET method",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed: GET\n",
		},
		{
			name:       "Invalid PUT method",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed: PUT\n",
		},
		{
			name:       "Invalid DELETE method",
			method:     http.MethodDelete,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed: DELETE\n",
		},
		{
			name:       "Invalid PATCH method",
			method:     http.MethodPatch,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed: PATCH\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/introspect", nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %v", err)
			}

			w := newMockResponseWriter()
			got := validateHTTPMethod(w, req)

			// Check return value
			if tt.method == http.MethodPost {
				if !got {
					t.Error("validateHTTPMethod() = false, want true for POST method")
				}
			} else {
				if got {
					t.Error("validateHTTPMethod() = true, want false for non-POST method")
				}
			}

			// Check status code
			if w.statusCode != tt.wantStatus {
				t.Errorf("validateHTTPMethod() status = %v, want %v", w.statusCode, tt.wantStatus)
			}

			// Check response body
			if string(w.body) != tt.wantBody {
				t.Errorf("validateHTTPMethod() body = %q, want %q", string(w.body), tt.wantBody)
			}
		})
	}
}
