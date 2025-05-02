# This Makefile is intended for development purposes only.
# It provides convenience commands for testing and debugging the key management functionality.

.PHONY: build run-generate run-list run-delete clean

# Build the CLI
build:
	go build -o bin/keys main.go

# Generate a new key pair
run-generate: build
	./bin/keys generate

# List available key pairs
run-list: build
	./bin/keys list

# Delete a key pair (requires KEY_ID)
run-delete: build
	@if [ -z "$(KEY_ID)" ]; then \
		echo "Error: KEY_ID is required. Usage: make run-delete KEY_ID=<key-id>"; \
		exit 1; \
	fi
	./bin/keys delete -id $(KEY_ID)

# Clean build artifacts
clean:
	rm -rf bin/
