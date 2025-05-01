# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- RSA key management tool in `pkg/keys` for JWT signing key generation
- CLI interface for key management operations (generate, list, delete)
- Test keys for development and testing purposes
- RS256 signing for JWT tokens
- Environment variable configuration for JWT signature key file
- Documentation for key management tool and usage
- User pool implementation for managing client credentials
- Documentation for user pool configuration and management
- Default test user credentials for development
- GitHub Actions workflow for basic branch builds
- Basic Authentication validation for client credentials
- Base64 encoding validation for credentials
- Proper OAuth2 error responses according to RFC 6749
- Test utilities for token endpoint testing
- `/token` endpoint with Basic Authentication support
- Initial implementation of OAuth2 server with client credentials grant
- Task documentation
- Initial project setup

### Changed
- Switched from HS256 to RS256 for JWT token signing
- Updated JWT token generation to use RSA private keys
- Extracted user management into dedicated `internal/userpool` package
- Updated README with key management and user pool configuration details
- Enhanced error logging using slog throughout the application
- Refactored authentication and token handling into internal packages:
  - Moved token generation to `internal/token` package
  - Moved Basic Auth handling to `internal/auth` package
- Improved error handling and response formatting
- Better separation of concerns and code organization

### Deleted
- HS256 signing implementation
- JWT secret key environment variable configuration
- Static JWT token response
- Direct key management in main application
