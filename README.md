# OAuth2 Server

A simple OAuth2 server implementation in Go that supports the Client Credentials Grant flow with Basic Authentication.

## Features

- OAuth2 Client Credentials Grant flow ([RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749))
- JWT Access Token issuance ([RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)) with RS256 signing
- Basic Authentication for client credentials
- Token introspection endpoint ([RFC 7662](https://datatracker.ietf.org/doc/html/rfc7662))
- JWK endpoint for signing keys ([RFC 7517](https://datatracker.ietf.org/doc/html/rfc7517))
- Local deployment using k3d (Kubernetes in Docker)

## Prerequisites

### Core Requirements
- [Go](https://golang.org/dl/) 1.21 or later (required for building and running the server)

### Development Tools
- [Make](https://www.gnu.org/software/make/) (optional, for using Makefile commands in keytool)

### Local Kubernetes Deployment
- [Docker](https://docs.docker.com/get-docker/) (required for building and running containers)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) (required for Kubernetes cluster interaction)
- [k3d](https://k3d.io/v5.6.0/#installation) (required for local Kubernetes cluster)

## Getting Started

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/wunderkind2k1/auth-server-task.git
cd auth-server-task
```

2. Generate an RSA key pair for JWT signing:
```bash
cd keytool
make run-generate
```

3. Start the server with the generated key:
```bash
# Export the private key content (replace <keyID> with your actual key ID)
export JWT_SIGNATURE_KEY="$(cat keytool/keys/<keyID>.private.pem)"
go run server/main.go
```

Alternatively, you can use the convenience script:
```bash
cd server
./startServer.sh
```

The server will start on port 8080.

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| JWT_SIGNATURE_KEY | Content of the RSA private key in PEM format for JWT signing | Yes |

### Key Management

The project includes a separate key management tool in the `keytool` directory. This tool provides commands for:
- Generating RSA key pairs
- Listing available keys
- Deleting key pairs

For development, you can use the pre-generated test keys in `keytool/keys/`. See [keytool/README.md](keytool/README.md) for more details.

### User Pool Configuration

The server uses a simple in-memory user pool for authentication. By default, it includes a test user with the following credentials:
- Client ID: `sho`
- Client Secret: `test123`

To modify the user pool, you can edit the [`server/internal/userpool/default.go`](server/internal/userpool/default.go) file. The user pool is implemented as a simple map structure where you can add, remove, or modify users. Each user requires a client ID and client secret.

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

### Local Deployment with k3d

For a more production-like environment, you can deploy the server using k3d:

1. Create a local Kubernetes cluster:
```bash
cd deployment/local
./manage-cluster.sh create
```

2. Set up the JWT secret:
```bash
./setup-secret.sh
```

3. Build and deploy the application:
```bash
./rebuildServerImage.sh
./deploy.sh
```

4. Verify the deployment:
```bash
./check-deployment.sh
```

The server will be accessible at `http://localhost:8080`. See [deployment/local/README.md](deployment/local/README.md) for detailed deployment instructions.

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

### JWKS Endpoint

Provides the JSON Web Key Set (JWKS) for token verification. The endpoint follows RFC 7517 and only accepts GET requests.

```bash
curl -X GET http://localhost:8080/.well-known/jwks.json
```

Response:
```json
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "1",
      "alg": "RS256",
      "n": "...",
      "e": "..."
    }
  ]
}
```

### Token Introspection Endpoint

Validates and provides information about an access token. The endpoint follows RFC 7662 and requires Basic Authentication.

```bash
curl -X POST http://localhost:8080/introspect \
  -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

Response for valid token:
```json
{
  "active": true,
  "scope": "",
  "client_id": "sho",
  "username": "sho",
  "token_type": "Bearer",
  "exp": 1735689600,
  "iat": 1735686000,
  "nbf": 1735686000,
  "sub": "sho",
  "aud": [],
  "iss": "https://oauth2-server",
  "jti": "unique-token-id"
}
```

Response for invalid token:
```json
{
  "active": false
}
```

## Testing

Test scripts are provided to verify the functionality of both endpoints:

```bash
# Set the JWT_SIGNATURE_KEY environment variable first
export JWT_SIGNATURE_KEY="$(cat keytool/keys/<keyID>.private.pem)"

# Test token endpoint
./test-utils/test-token-endpoint.sh

# Test JWKS endpoint
./test-utils/test-jwks-endpoint.sh

# Test introspection endpoint
./test-utils/test-introspection-endpoint.sh
```

## Development

The project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html) and maintains a [CHANGELOG.md](CHANGELOG.md) following the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format.

### Enforcing Note on Code Duplication

As Rob Pike famously said in his talk ["Go Proverbs"](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s), "A little copying is better than a little dependency." This principle is applied in our codebase, particularly in the key management implementation, where we intentionally maintain some code duplication between the main application and the key management tool to avoid tight coupling.

The key parsing logic in the main application ([`internal/token/keyloader.go`](../server/internal/token/keyloader.go)) intentionally duplicates some code from this package. This is by design to maintain separation between the core application and its supporting tools. The main application should be self-contained and not depend on this package.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
