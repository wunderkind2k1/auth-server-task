package userpool

// Default returns a user pool with default test users
func Default() map[string]string {
	return map[string]string{
		"sho": "test123",
	}
}
