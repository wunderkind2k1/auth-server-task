package request

import (
	"net/http"
	"testing"
)

// mockResponseWriter is a simple mock of http.ResponseWriter
type mockResponseWriter struct {
	statusCode int
	header     http.Header
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		header: make(http.Header),
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.header
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedMethod string
		wantValid      bool
		wantStatus     int
	}{
		{
			name:           "valid method",
			method:         http.MethodPost,
			expectedMethod: http.MethodPost,
			wantValid:      true,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			expectedMethod: http.MethodPost,
			wantValid:      false,
			wantStatus:     http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Method: tt.method}
			w := newMockResponseWriter()
			got := ValidateMethod(w, req, tt.expectedMethod)
			if got != tt.wantValid {
				t.Errorf("ValidateMethod() = %v, want %v", got, tt.wantValid)
			}
			if !tt.wantValid && w.statusCode != tt.wantStatus {
				t.Errorf("ValidateMethod() status = %v, want %v", w.statusCode, tt.wantStatus)
			}
		})
	}
}

func TestValidateContentType(t *testing.T) {
	tests := []struct {
		name                string
		contentType         string
		expectedContentType string
		wantValid           bool
		wantStatus          int
	}{
		{
			name:                "valid content type",
			contentType:         "application/json",
			expectedContentType: "application/json",
			wantValid:           true,
		},
		{
			name:                "invalid content type",
			contentType:         "text/plain",
			expectedContentType: "application/json",
			wantValid:           false,
			wantStatus:          http.StatusBadRequest,
		},
		{
			name:                "missing content type",
			contentType:         "",
			expectedContentType: "application/json",
			wantValid:           false,
			wantStatus:          http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: make(http.Header),
			}
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			w := newMockResponseWriter()
			got := ValidateContentType(w, req, tt.expectedContentType)
			if got != tt.wantValid {
				t.Errorf("ValidateContentType() = %v, want %v", got, tt.wantValid)
			}
			if !tt.wantValid && w.statusCode != tt.wantStatus {
				t.Errorf("ValidateContentType() status = %v, want %v", w.statusCode, tt.wantStatus)
			}
		})
	}
}

func TestValidateAuthorization(t *testing.T) {
	tests := []struct {
		name           string
		auth           string
		expectedScheme string
		wantValid      bool
		wantStatus     int
	}{
		{
			name:           "valid authorization",
			auth:           "Bearer token123",
			expectedScheme: "Bearer",
			wantValid:      true,
		},
		{
			name:           "invalid scheme",
			auth:           "Basic token123",
			expectedScheme: "Bearer",
			wantValid:      false,
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name:           "missing authorization",
			auth:           "",
			expectedScheme: "Bearer",
			wantValid:      false,
			wantStatus:     http.StatusUnauthorized,
		},
		{
			name:           "malformed authorization",
			auth:           "Bearer",
			expectedScheme: "Bearer",
			wantValid:      false,
			wantStatus:     http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: make(http.Header),
			}
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			w := newMockResponseWriter()
			got := ValidateAuthorization(w, req, tt.expectedScheme)
			if got != tt.wantValid {
				t.Errorf("ValidateAuthorization() = %v, want %v", got, tt.wantValid)
			}
			if !tt.wantValid && w.statusCode != tt.wantStatus {
				t.Errorf("ValidateAuthorization() status = %v, want %v", w.statusCode, tt.wantStatus)
			}
		})
	}
}
