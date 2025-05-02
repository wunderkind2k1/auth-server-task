# Key Management Tool

This package provides a standalone CLI tool for managing RSA key pairs used for JWT signing in the main application. It is designed as a separate utility that supports the main application but operates independently.

## Purpose

The key management tool serves several important purposes:
- Generates RSA key pairs for JWT signing
- Manages key storage in a structured way
- Provides a simple CLI interface for key operations
- Supports the main application's JWT signing requirements

## Package Structure

This package is intentionally separated from the main application code to maintain a clear boundary between the core application logic and its supporting tools. While it exists in the same repository (mono repo approach), it should not be imported or used directly by the main application.

The main application uses its own key loading logic in [`internal/token/keyloader.go`](../server/internal/token/keyloader.go) to maintain independence and avoid coupling with this tool.

## Test Keys

The [`keys/`](keys/) directory contains pre-generated key pairs for development and testing purposes. See [keys/README.md](keys/README.md) for details about these test keys and their usage.

## CLI Tool (`internal/main.go`)

The CLI tool provides the following commands:

### Generate Key Pair
```bash
./bin/keys generate
```
Generates a new RSA key pair with default 2048-bit key size. The keys are saved in the `keys` directory with the following format:
- Private key: `<keyID>.private.pem`
- Public key: `<keyID>.public.pem`

### List Key Pairs
```bash
./bin/keys list
```
Lists all available key pairs in the `keys` directory, showing their IDs and file paths.

### Delete Key Pair
```bash
./bin/keys delete -id <keyID>
```
Deletes a specific key pair by its ID.

## Makefile

The Makefile provides development convenience commands for working with the key management tool. It is intended for development purposes only and should not be used in production.

Available commands:
- `make build`: Builds the CLI tool
- `make run-generate`: Generates a new key pair
- `make run-list`: Lists available key pairs
- `make run-delete KEY_ID=<keyID>`: Deletes a specific key pair
- `make clean`: Removes build artifacts

## Usage in Main Application

The main application expects a private key file to be available at the path specified by the `JWT_SIGNATURE_KEY_FILE` environment variable. This file should be generated using this tool:

1. Generate a key pair:
```bash
cd keytool
make run-generate
```

2. Use the generated private key in the main application:
```bash
JWT_SIGNATURE_KEY_FILE=keytool/keys/<keyID>.private.pem go run ../server/main.go
```

For development, you can also use the pre-generated test keys in the [`keys/`](keys/) directory.

## Enforcing Note on Code Duplication

The key loading logic in the main application ([`internal/token/keyloader.go`](../server/internal/token/keyloader.go)) intentionally duplicates some code from this package. This is by design to maintain separation between the core application and its supporting tools. The main application should be self-contained and not depend on this package.
