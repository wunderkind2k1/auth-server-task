# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.0.10] - 2025-05-07

### Added
- Added test ratio calculation and enforcement:
  - Minimum threshold of 0.5 (50% test-to-code)
  - Warning threshold of 0.7 (70% test-to-code)
  - Target threshold of 1.0 (100% test-to-code)
  - Component-specific ratios for keytool and server
  - HTML coverage reports as artifacts
- Added comprehensive RSA key pair validation using crypto/rsa.Validate()
- Added detailed package documentation for keytool/internal/rsa package
- Added test coverage for RSA key pair validation
- Added request validation package with HTTP method, content type, and authorization header validation
- Added comprehensive test coverage for request validation
- Added nolint directive for TestGenerateToken to maintain test organization
- Initial test coverage for token introspection
  - HTTP method validation (POST-only as per RFC 7662)
  - Token extraction from request
  - Token validation and signing method checks
  - Error response handling
- Initial test coverage for basic auth functionality
- Extracted HTTP method validation into a dedicated function
- Improved error response handling for token validation failures
- Added linting merge gate to GitHub Actions workflow:
  - Configured to run on all pushes and PRs
  - Acts as a merge gate for main branch
  - Allows feature branch pushes even with linting issues
- Added GitHub Actions workflow documentation:
  - Created `.github/workflows/README.md` with detailed workflow explanations
  - Added CI/CD section to main README.md
- Added deployment management scripts:
  - `scale.sh` for scaling the OAuth2 server deployment
  - `undeploy.sh` for removing the deployment
- Added golangci-lint configuration for code quality
- Added Makefile with common development commands
- Added comprehensive token introspection endpoint tests
- Added proper error handling for token introspection according to RFC 7662
- Added reference to Rob Pike's Go proverb about code duplication in the context of key management implementation

### Changed
- Improved GitHub Actions workflow:
  - Moved workflow documentation to `.github/README.md`
  - Removed duplicate documentation from workflows directory
  - Updated main README to reference new documentation location
- Improved RSA key pair validation with built-in crypto/rsa.Validate()
- Enhanced keytool package documentation with security considerations and usage examples
- Simplified test structure by removing unnecessary setup functions
- Refactored request validation to use custom mock ResponseWriter instead of httptest
- Improved error handling in JWKS response writing
- Simplified embedded field selectors in keyloader tests
- Updated Go version requirement to 1.24
- Improved key file naming consistency in RSA key management
- Enhanced error handling in key management operations
- Centralized response handling in token introspection
- Improved documentation with proper package comments
- Enhanced security with proper file permissions for key files
- Added proper error handling for JSON encoding failures
- Improved server configuration with read header timeout
- Improved README organization and clarity:
  - Reordered sections for better flow
  - Categorized prerequisites by purpose
  - Added links to prerequisite tools
  - Added comprehensive token introspection endpoint documentation
  - Enhanced API endpoint documentation with examples

### Fixed
- Fixed error handling in JWKS response writing
- Fixed unused parameters in token validation tests
- Fixed missing periods in documentation comments
- Fixed token introspection response for invalid tokens (now returns 200 with active=false)
- Fixed key file naming inconsistency in RSA key management
- Fixed error handling in token introspection endpoint
- Fixed potential security issues with file permissions
- Fixed error handling in JSON response encoding

## [0.0.9] - 2024-05-03

### Added
- Local deployment setup using k3d (Kubernetes in Docker)
- Deployment scripts for local development:
  - `manage-cluster.sh` for k3d cluster management
  - `deploy.sh` for application deployment
  - `check-deployment.sh` for deployment verification
  - `setup-secret.sh` for JWT key management
  - `rebuildServerImage.sh` for Docker image building
- Kubernetes manifests for local deployment
- Documentation for local deployment setup and usage
- Token introspection endpoint implementation according to RFC 7662
- Test utilities for token introspection endpoint
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
- Improved local development experience with k3d-based deployment
- Enhanced deployment verification with comprehensive checks
- Updated documentation to include local deployment instructions
- Improved security by using environment variable for JWT signing key content instead of file path
- Updated package documentation to reflect token introspection functionality
- Enhanced key parsing logic to support environment variable-based configuration
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
- Suppressed Docker images in git to reduce repository size

### Deleted
- Old `pkg/keys` directory structure
- Obsolete userpool.go file
- Direct key management in main application
- Static JWT token response
- JWT secret key environment variable configuration
- HS256 signing implementation
