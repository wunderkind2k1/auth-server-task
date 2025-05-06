// Package request provides HTTP request validation utilities for OAuth2 endpoints.
// It implements common validation requirements for OAuth2 endpoints as specified in RFC 7662,
// including method validation, content type validation, and authorization header validation.
//
// The package is designed to be used in HTTP handlers to validate incoming requests
// before processing them. It provides functions to validate:
//   - HTTP method (e.g., ensuring POST for introspection endpoint)
//   - Content-Type header (e.g., ensuring application/json)
//   - Authorization header (e.g., validating Bearer token format)
//
// Each validation function writes appropriate error responses to the http.ResponseWriter
// when validation fails, following OAuth2 error response format.
package request

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// ValidateMethod checks if the request method matches the expected method.
// If the method doesn't match, it writes a MethodNotAllowed error response and returns false.
// Returns true if the method is valid.
func ValidateMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		slog.Error("Method not allowed", "method", r.Method, "expected", expectedMethod, "status", http.StatusMethodNotAllowed)
		http.Error(w, fmt.Sprintf("Method not allowed: %s", r.Method), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// ValidateContentType checks if the request has the expected Content-Type header.
// If the Content-Type doesn't match, it writes a BadRequest error response and returns false.
// Returns true if the Content-Type is valid.
func ValidateContentType(w http.ResponseWriter, r *http.Request, expectedContentType string) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != expectedContentType {
		slog.Error("Invalid Content-Type", "got", contentType, "expected", expectedContentType, "status", http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Invalid Content-Type: %s", contentType), http.StatusBadRequest)
		return false
	}
	return true
}

// ValidateAuthorization checks if the request has a valid Authorization header with the expected scheme.
// If the Authorization header is missing or invalid, it writes an Unauthorized error response and returns false.
// Returns true if the Authorization header is valid.
func ValidateAuthorization(w http.ResponseWriter, r *http.Request, expectedScheme string) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		slog.Error("Authorization header required", "status", http.StatusUnauthorized)
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return false
	}

	if !strings.HasPrefix(auth, expectedScheme+" ") {
		slog.Error("Invalid authorization scheme", "got", auth, "expected", expectedScheme, "status", http.StatusUnauthorized)
		http.Error(w, fmt.Sprintf("Invalid authorization scheme: %s", auth), http.StatusUnauthorized)
		return false
	}

	return true
}
