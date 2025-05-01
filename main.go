package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"oauth2-task/internal/auth"
	"oauth2-task/internal/token"
	"oauth2-task/internal/userpool"
)

var (
	secretKey string
	userPool  map[string]string
)

func setup() {
	secretKey = os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		slog.Error("Mandatory JWT_SECRET_KEY environment variable is not set")
		os.Exit(1)
	}
	slog.Info("Mandatory JWT_SECRET_KEY environment variable is set")

	// Initialize user pool with default test users
	userPool = userpool.Default()
}

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

	// Create BasicAuth instance with user pool
	basicAuth := auth.NewBasicAuth(userPool)

	// Validate Basic Auth
	if err := basicAuth.ParseBasicAuth(r.Header.Get("Authorization")); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := auth.GetErrorResponse(err)
		json.NewEncoder(w).Encode(errorResponse)
		slog.Error("Authentication failed", "error", err)
		return
	}

	// Create token generator
	generator := token.NewGenerator(secretKey)

	// Generate a real JWT token
	tokenString, err := generator.GenerateToken()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(auth.ErrorResponse{
			Error:            "server_error",
			ErrorDescription: "Failed to generate token",
		})
		if err != nil {
			slog.Error("Failed to generate Error Response for token generation", "error", err)
			return
		}
		slog.Error("Failed to generate token", "error", err)
		return
	}

	// Return the token response
	response := TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	setup()
	server := &http.Server{Addr: ":8080"}
	slog.Info("Starting server", "port", 8080)
	http.HandleFunc("/token", tokenHandler)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
