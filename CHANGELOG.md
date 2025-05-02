# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Convenience script `startServer.sh` for quick server startup during development
- RSA key management tool in `keytool` for JWT signing key generation
- Package documentation for all internal packages
- Test utilities for JWKS endpoint testing
- JWKS endpoint implementation according to RFC 7517
- Initial project setup
- Task documentation
- Initial implementation of OAuth2 server with client credentials grant
- `/token` endpoint with Basic Authentication support
- Test utilities for token endpoint testing
- Proper OAuth2 error responses according to RFC 6749
- Base64 encoding validation for credentials
- Basic Authentication validation for client credentials
- GitHub Actions workflow for basic branch builds
- Default test user credentials for development
- Documentation for user pool configuration and management
- User pool implementation for managing client credentials
- Documentation for key management tool and usage
- Environment variable configuration for JWT signature key file
- RS256 signing for JWT tokens
- Test keys for development and testing purposes
- CLI interface for key management operations (generate, list, delete)

### Changed
- Improved security by using environment variable for JWT signing key content instead of file path
- Restructured project into a monorepo with separate `keytool` and `server` components
- Updated GitHub Actions workflow to handle multiple Go modules
- Updated all documentation to reflect new repository structure
- Added method validation for JWKS endpoint
- Cleaned up userpool package structure
- Better separation of concerns and code organization
- Improved error handling and response formatting
- Refactored authentication and token handling into internal packages:
  - Moved token generation to `internal/token` package
  - Moved Basic Auth handling to `internal/auth` package
- Enhanced error logging using slog throughout the application
- Updated README with key management and user pool configuration details
- Extracted user management into dedicated `internal/userpool` package
- Updated JWT token generation to use RSA private keys
- Switched from HS256 to RS256 for JWT token signing

### Deleted
- Old `pkg/keys` directory structure
- Obsolete userpool.go file
- Direct key management in main application
- Static JWT token response
- JWT secret key environment variable configuration
- HS256 signing implementation
