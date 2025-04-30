# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions workflow for basic branch builds
- Environment variable configuration for JWT secret key
- JWT token generation using golang-jwt library
- Basic Authentication validation for client credentials
- Base64 encoding validation for credentials
- Proper OAuth2 error responses according to RFC 6749
- Test utilities for token endpoint testing
- Static JWT token response
- `/token` endpoint with Basic Authentication support
- Initial implementation of OAuth2 server with client credentials grant
- Task documentation
- Initial project setup

### Changed
  - Enhanced error logging using slog throughout the application
  - Moved JWT secret key to environment variable configuration
  - Refactored authentication and token handling into internal packages:
  - Moved token generation to `internal/token` package
  - Moved Basic Auth handling to `internal/auth` package
  - Improved error handling and response formatting
  - Better separation of concerns and code organization
