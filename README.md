# OAuth2 Server

A simple OAuth2 server implementation in Go that supports the Client Credentials Grant flow with Basic Authentication.

## Features

- OAuth2 Client Credentials Grant flow ([RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749))
- JWT Access Token issuance ([RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)) with RS256 signing
- Basic Authentication for client credentials
- Token introspection endpoint ([RFC 7662](https://datatracker.ietf.org/doc/html/rfc7662))
- JWK endpoint for signing keys ([RFC 7517](https://datatracker.ietf.org/doc/html/rfc7517))

## Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/wunderkind2k1/auth-server-task.git
cd auth-server-task
```

2. Generate an RSA key pair for JWT signing:
```bash
cd pkg/keys
make run-generate
```

3. Start the server with the generated key:
```bash
JWT_SIGNATURE_KEY_FILE=pkg/keys/keys/<keyID>.private.pem go run main.go
```

The server will start on port 8080.

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| JWT_SIGNATURE_KEY_FILE | Path to the RSA private key file for JWT signing | Yes |

### Key Management

The project includes a separate key management tool in the `pkg/keys` directory. This tool provides commands for:
- Generating RSA key pairs
- Listing available keys
- Deleting key pairs

For development, you can use the pre-generated test keys in `pkg/keys/keys/`. See [pkg/keys/README.md](pkg/keys/README.md) for more details.

### User Pool Configuration

The server uses a simple in-memory user pool for authentication. By default, it includes a test user with the following credentials:
- Client ID: `sho`
- Client Secret: `test123`

To modify the user pool, you can edit the [`internal/userpool/default.go`](internal/userpool/default.go) file. The user pool is implemented as a simple map structure where you can add, remove, or modify users. Each user requires a client ID and client secret.

Example of adding a new user:
```go
func Default() map[string]string {
    return map[string]string{
        "sho": "test123",
        "new-user": "new-secret",
    }
}
```

Note: In a production environment, you should implement a more secure and persistent storage solution for user credentials.

## API Endpoints

### Token Endpoint

Issues JWT access tokens using the Client Credentials Grant flow. Tokens are signed using RS256.

```bash
curl -X POST http://localhost:8080/token \
  -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
  -H "Content-Type: application/json"
```

Response:
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Testing

A test script is provided to verify the token endpoint functionality:

```bash
./test-utils/test-token-endpoint.sh
```

## Development

The project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html) and maintains a [CHANGELOG.md](CHANGELOG.md) following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
