package main

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// TokenError represents an OAuth2 error response
type TokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		slog.Error("Method not allowed", "method", r.Method, "status", http.StatusMethodNotAllowed)
		return
	}

	// Check Basic Auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TokenError{
			Error:            "invalid_client",
			ErrorDescription: "Authorization header required",
		})
		slog.Error("Missing authorization header")
		return
	}

	// Parse Basic Auth
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TokenError{
			Error:            "invalid_client",
			ErrorDescription: "Invalid authorization header format",
		})
		slog.Error("Invalid authorization header format")
		return
	}

	// For now, we'll accept any valid Basic Auth credentials
	// In a real implementation, we would validate against stored client credentials
	_, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TokenError{
			Error:            "invalid_client",
			ErrorDescription: "Invalid credentials",
		})
		slog.Error("Invalid base64 encoding in credentials")
		return
	}

	// Return a static token response
	response := TokenResponse{
		AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNTE2MjM5MDIyfQ",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	server := &http.Server{Addr: ":8080"}
	slog.Info("Starting server", "port", 8080)
	http.HandleFunc("/token", tokenHandler)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
