package auth

import (
	"encoding/base64"
	"testing"
)

func TestParseBasicAuth(t *testing.T) {
	// Setup test user pool
	pool := map[string]string{
		"testuser": "testpass",
		"admin":    "adminpass",
	}
	ba := NewBasicAuth(pool)

	tests := []struct {
		name         string
		authHeader   string
		wantErr      error
		wantUsername string
		wantPassword string
	}{
		{
			name:         "Valid credentials",
			authHeader:   "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass")),
			wantErr:      nil,
			wantUsername: "testuser",
			wantPassword: "testpass",
		},
		{
			name:       "Missing header",
			authHeader: "",
			wantErr:    ErrMissingHeader,
		},
		{
			name:       "Invalid format - no space",
			authHeader: "Basic" + base64.StdEncoding.EncodeToString([]byte("testuser:testpass")),
			wantErr:    ErrInvalidFormat,
		},
		{
			name:       "Invalid format - wrong scheme",
			authHeader: "Bearer " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass")),
			wantErr:    ErrInvalidFormat,
		},
		{
			name:       "Invalid base64",
			authHeader: "Basic invalid-base64",
			wantErr:    ErrInvalidBase64,
		},
		{
			name:       "Invalid credentials format - no colon",
			authHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser")),
			wantErr:    ErrInvalidFormat,
		},
		{
			name:       "Invalid credentials - wrong password",
			authHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:wrongpass")),
			wantErr:    ErrInvalidCredentials,
		},
		{
			name:       "Invalid credentials - non-existent user",
			authHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("nonexistent:pass")),
			wantErr:    ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ba.ParseBasicAuth(tt.authHeader)

			// Check error
			if err != tt.wantErr {
				t.Errorf("ParseBasicAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected no error, check the credentials
			if err == nil {
				if ba.Username != tt.wantUsername {
					t.Errorf("ParseBasicAuth() username = %v, want %v", ba.Username, tt.wantUsername)
				}
				if ba.Password != tt.wantPassword {
					t.Errorf("ParseBasicAuth() password = %v, want %v", ba.Password, tt.wantPassword)
				}
			}
		})
	}
}

func TestGetErrorResponse(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantError string
		wantDesc  string
	}{
		{
			name:      "Missing header error",
			err:       ErrMissingHeader,
			wantError: "invalid_client",
			wantDesc:  "Authorization header required",
		},
		{
			name:      "Invalid format error",
			err:       ErrInvalidFormat,
			wantError: "invalid_client",
			wantDesc:  "Invalid authorization header format",
		},
		{
			name:      "Invalid base64 error",
			err:       ErrInvalidBase64,
			wantError: "invalid_client",
			wantDesc:  "Invalid credentials",
		},
		{
			name:      "Invalid credentials error",
			err:       ErrInvalidCredentials,
			wantError: "invalid_client",
			wantDesc:  "Invalid username or password",
		},
		{
			name:      "Unknown error",
			err:       ErrInvalidAuthScheme,
			wantError: "invalid_client",
			wantDesc:  "Authentication failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorResponse(tt.err)
			if got.Error != tt.wantError {
				t.Errorf("GetErrorResponse() error = %v, want %v", got.Error, tt.wantError)
			}
			if got.ErrorDescription != tt.wantDesc {
				t.Errorf("GetErrorResponse() description = %v, want %v", got.ErrorDescription, tt.wantDesc)
			}
		})
	}
}
