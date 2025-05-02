// Package userpool manages user authentication credentials.
// It provides functionality for storing and validating user credentials
// in a simple in-memory user pool. The package is used for Basic
// Authentication in the OAuth2 token endpoint.
package userpool

// Default returns a user pool with default test users
func Default() map[string]string {
	return map[string]string{
		"sho": "test123",
	}
}
