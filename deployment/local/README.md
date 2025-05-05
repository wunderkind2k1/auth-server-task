# Local Deployment

This directory contains scripts and configurations for deploying the OAuth2 server locally using Docker and Kubernetes (k3d).

## Prerequisites

- Docker
- kubectl
- k3d (Lightweight Kubernetes distribution)

## Available Scripts

### Cluster Management
- `manage-cluster.sh`: Manages the k3d cluster
  - `./manage-cluster.sh create`: Creates a new cluster with port 8080 mapped to the loadbalancer
  - `./manage-cluster.sh delete`: Deletes the cluster
  - `./manage-cluster.sh status`: Checks if the cluster is running

### Deployment
- `deploy.sh`: Deploys the OAuth2 server to the cluster
  - Loads the Docker image into the cluster
  - Applies Kubernetes manifests
  - Creates necessary resources

- `undeploy.sh`: Removes the OAuth2 server deployment
  - Deletes the deployment from the cluster
  - Useful for clean removal without deleting the cluster

- `scale.sh`: Scales the OAuth2 server deployment
  - Usage: `REPLICAS=2 ./scale.sh`
  - Adjusts the number of running pods
  - Requires REPLICAS environment variable

### Image Management
- `rebuildServerImage.sh`: Builds and saves the Docker image
  - Builds the server image
  - Saves it as a tar file for local use
  - Run this after code changes

### Configuration
- `setup-secret.sh`: Sets up the JWT signing key
  - Creates a Kubernetes secret from the first private key in keytool
  - Only needed for initial setup or key changes

### Verification
- `check-deployment.sh`: Verifies the deployment
  - Checks pod status and readiness
  - Verifies service availability
  - Tests loadbalancer and port-forwarding
  - Validates JWKS endpoint

- `port-forward.sh`: Sets up port forwarding
  - Forwards local port 8080 to the service
  - Alternative to loadbalancer access

## Development Flow

1. **Initial Setup** (one-time):
   ```bash
   ./manage-cluster.sh create
   ./setup-secret.sh
   ```

2. **Development Cycle**:
   ```bash
   # After code changes
   ./rebuildServerImage.sh
   ./deploy.sh
   ./check-deployment.sh
   ```

3. **Scaling** (if needed):
   ```bash
   REPLICAS=2 ./scale.sh
   ```

4. **Cleanup**:
   ```bash
   # Remove deployment
   ./undeploy.sh

   # Or remove entire cluster
   ./manage-cluster.sh delete
   ```

## Security Notes

- The JWT signing key is managed as a Kubernetes secret, not embedded in the container image
- The secret persists in the cluster until explicitly deleted

## Service Access

The OAuth2 server is exposed as a LoadBalancer service:
- Access via k3d loadbalancer: `http://localhost:8080`
- Access via port-forwarding: `kubectl port-forward svc/oauth2-server 8080:8080`
- JWKS endpoint: `http://localhost:8080/.well-known/jwks.json`

## Troubleshooting

If you encounter issues:
1. Check cluster status: `k3d cluster list`
2. Verify secret: `kubectl get secret jwt-key`
3. Check pod logs: `kubectl logs -l app=oauth2-server`
4. Check pod status: `kubectl get pods -l app=oauth2-server`
5. Check service status: `kubectl get svc oauth2-server`
6. Run deployment check: `./check-deployment.sh`

Common issues:
- Service not accessible: Check k3d loadbalancer with `docker ps | grep k3d`
- 404 errors: Verify pod and service configuration
- Cluster recreation: Run `./manage-cluster.sh delete` followed by `./manage-cluster.sh create`
