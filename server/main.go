package main

import (
	"encoding/pem"
	"log/slog"
	"net/http"
	"oauth2-task/internal/auth"
	"oauth2-task/internal/token"
	"oauth2-task/internal/userpool"
	"os"
)

var (
	keyPair  token.KeyPair
	userPool map[string]string
)

func setup() {
	// Get key content from environment variable
	keyContent := os.Getenv("JWT_SIGNATURE_KEY")
	if keyContent == "" {
		slog.Error("Mandatory JWT_SIGNATURE_KEY environment variable is not set")
		os.Exit(1)
	}

	// Parse the PEM block
	block, _ := pem.Decode([]byte(keyContent))
	if block == nil {
		slog.Error("Failed to decode PEM block from JWT_SIGNATURE_KEY")
		os.Exit(1)
	}

	// Load the private key
	var err error
	keyPair, err = token.ParsePrivateKey(block.Bytes)
	if err != nil {
		slog.Error("Failed to parse private key", "error", err)
		os.Exit(1)
	}
	slog.Info("Private key loaded successfully from environment variable")

	// Initialize user pool with default test users
	userPool = userpool.Default()
}

func main() {
	setup()
	server := &http.Server{Addr: ":8080"}
	slog.Info("Starting server", "port", 8080)
	http.HandleFunc("/token", auth.HandleToken(keyPair, userPool))
	http.HandleFunc("/.well-known/jwks.json", auth.HandleJWKS(keyPair))
	http.HandleFunc("/introspect", token.HandleIntrospection(keyPair))
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
