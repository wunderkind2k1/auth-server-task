// Package userpool provides client credential management for OAuth2 authentication.
// It implements a simple in-memory storage for client credentials (client_id/client_secret pairs)
// used in the OAuth2 Client Credentials Grant flow. The package is designed to be
// easily extensible for different storage backends in production environments.
package userpool

// Default returns a user pool with default test users.
// This function is intended for development and testing purposes only.
// In production, implement a proper credential storage solution.
func Default() map[string]string {
	return map[string]string{
		"sho": "test123",
	}
}
