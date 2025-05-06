package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"oauth2-task/internal/request"
	"oauth2-task/internal/token"
)

// TokenResponse represents the OAuth2 token response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// HandleToken processes OAuth2 token requests.
func HandleToken(keyPair token.KeyPair, userPool map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !request.ValidateMethod(w, r, http.MethodPost) {
			return
		}

		// Create BasicAuth instance with user pool
		basicAuth := NewBasicAuth(userPool)

		// Validate Basic Auth
		if !request.ValidateAuthorization(w, r, "Basic") {
			return
		}

		// Parse Basic Auth credentials
		if err := basicAuth.ParseBasicAuth(r.Header.Get("Authorization")); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := GetErrorResponse(err)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				slog.Error("Failed to encode error response", "error", err)
				return
			}
			slog.Error("Authentication failed", "error", err)
			return
		}

		// Create token generator
		generator := token.NewGenerator(keyPair.PrivateKey())

		// Generate a real JWT token
		tokenString, err := generator.GenerateToken(basicAuth.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{
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
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Failed to return token", "error", err)
			return
		}
	}
}
