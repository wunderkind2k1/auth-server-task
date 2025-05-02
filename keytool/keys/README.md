# Test Keys Directory

This directory contains RSA key pairs that are provided for development and testing purposes only. These keys are:

- Convenient for local development
- Ready to use for testing the JWT signing functionality
- Safe to delete and regenerate at any time
- **NOT** intended for production use

## Available Test Keys

The keys in this directory follow the standard naming convention:
- `<keyID>.private.pem`: Private key for JWT signing
- `<keyID>.public.pem`: Corresponding public key

## Usage

To use these test keys with the main application during development:

```bash
JWT_SIGNATURE_KEY_FILE=keytool/keys/<keyID>.private.pem go run ../server/main.go
```

## Important Notes

1. These keys are committed to the repository for development convenience only
2. They should never be used in production environments
3. You can safely delete and regenerate these keys using the [key management CLI](../README.md#cli-tool-internalmaingo)
4. For production, always generate new key pairs using the key management tool

## Security

While these keys are safe to commit to the repository (as they are only for development), it's good practice to:
- Never commit production keys
- Regularly rotate test keys
- Use different keys for different environments
