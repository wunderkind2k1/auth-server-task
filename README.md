# OAuth2 Server

A simple OAuth2 server implementation in Go that supports the Client Credentials Grant flow with Basic Authentication.

## Features

- OAuth2 Client Credentials Grant flow (RFC 6749)
- JWT Access Token issuance (RFC 7519)
- Basic Authentication for client credentials
- Token introspection endpoint (RFC 7662)
- JWK endpoint for signing keys (RFC 7517)

## Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/wunderkind2k1/auth-server-task.git
cd auth-server-task
```

2. Run the server:
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Token Endpoint

Issues JWT access tokens using the Client Credentials Grant flow.

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
