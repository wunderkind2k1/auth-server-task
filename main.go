package main

import (
	"crypto/rsa"
	"log/slog"
	"net/http"
	"os"

	"oauth2-task/internal/auth"
	"oauth2-task/internal/token"
	"oauth2-task/internal/userpool"
)

var (
	keyPair  *rsa.PrivateKey
	userPool map[string]string
)

func setup() {
	// Get key file path from environment variable
	keyFilePath := os.Getenv("JWT_SIGNATURE_KEY_FILE")
	if keyFilePath == "" {
		slog.Error("Mandatory JWT_SIGNATURE_KEY_FILE environment variable is not set")
		os.Exit(1)
	}

	// Ensure the key file exists
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		slog.Error("Mandatory JWT signature key file does not exist", "path", keyFilePath)
		os.Exit(1)
	}

	// Load the private key
	var err error
	keyPair, err = token.LoadPrivateKey(keyFilePath)
	if err != nil {
		slog.Error("Failed to load private key", "error", err, "path", keyFilePath)
		os.Exit(1)
	}
	slog.Info("Private key loaded successfully", "path", keyFilePath)

	// Initialize user pool with default test users
	userPool = userpool.Default()
}

func main() {
	setup()
	server := &http.Server{Addr: ":8080"}
	slog.Info("Starting server", "port", 8080)
	http.HandleFunc("/token", auth.HandleToken(keyPair, userPool))
	http.HandleFunc("/.well-known/jwks.json", auth.HandleJWKS)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
