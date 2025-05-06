package userpool

import (
	"testing"
)

func TestDefault(t *testing.T) {
	t.Run("returns expected test user", func(t *testing.T) {
		pool := Default()
		if pool == nil {
			t.Fatal("Expected non-nil user pool")
		}

		password, exists := pool["sho"]
		if !exists {
			t.Fatal("Default test user 'sho' not found")
		}
		if password != "test123" {
			t.Errorf("Expected password 'test123', got '%s'", password)
		}
	})

	t.Run("returns maps with same content", func(t *testing.T) {
		pool1 := Default()
		pool2 := Default()

		if len(pool1) != len(pool2) {
			t.Errorf("Expected same map size, got %d and %d", len(pool1), len(pool2))
		}
		if pool1["sho"] != pool2["sho"] {
			t.Error("Expected same content in both maps")
		}
	})

	t.Run("returns independent map instances", func(t *testing.T) {
		pool1 := Default()
		// Modify pool1 before getting pool2
		pool1["sho"] = "modified"

		pool2 := Default()
		if pool2["sho"] == "modified" {
			t.Error("Expected Default() to return a new map instance")
		}
		if pool2["sho"] != "test123" {
			t.Error("Expected Default() to return map with original content")
		}
	})
}
