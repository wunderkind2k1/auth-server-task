# Local Deployment

This directory contains scripts and configurations for deploying the OAuth2 server locally using Docker and Kubernetes (k3d).

## Prerequisites

- Docker
- kubectl
- k3d (Lightweight Kubernetes distribution)

## Preparations (One-time Setup)

These steps only need to be performed once when setting up the development environment:

1. Create a local Kubernetes cluster:
   ```bash
   ./manage-cluster.sh create
   ```
   This creates a k3d cluster with port 8080 mapped to the loadbalancer.

2. Set up the JWT secret:
   ```bash
   ./setup-secret.sh
   ```
   This script creates a Kubernetes secret from the first private key found in the keytool directory.

   Note: Only run this again if:
   - Creating a new cluster
   - Updating the secret with a new key
   - If the secret was accidentally deleted

## Development Flow

These are the steps you'll repeat during development:

1. Build and save the Docker image (after code changes):
   ```bash
   ./rebuildServerImage.sh
   ```
   This script builds the Docker image and saves it as a tar file for local use.

2. Deploy the application:
   ```bash
   ./deploy.sh
   ```
   This script loads the saved image into the cluster and applies the Kubernetes manifests.

3. Verify the deployment:
   ```bash
   ./check-deployment.sh
   ```
   This script checks:
   - Pod status and readiness
   - Service availability
   - Service accessibility through the k3d loadbalancer
   - Service accessibility via port-forwarding
   - JWKS endpoint availability

## Security Notes

- The JWT signing key is managed as a Kubernetes secret, not embedded in the container image
- The secret persists in the cluster until explicitly deleted

## Service Access

The OAuth2 server is exposed as a LoadBalancer service, which means:
- It's accessible through the k3d loadbalancer at `http://localhost:8080`
- The service is also accessible via port-forwarding:
  ```bash
  kubectl port-forward svc/oauth2-server 8080:8080
  ```
- The JWKS endpoint is available at `http://localhost:8080/.well-known/jwks.json`

## Cleanup

To remove the deployment and cluster:
```bash
./manage-cluster.sh delete
```

## Troubleshooting

If you encounter issues:
1. Check that the cluster is running: `k3d cluster list`
2. Verify the secret exists: `kubectl get secret jwt-key`
3. Check pod logs: `kubectl logs -l app=oauth2-server`
4. Check pod status: `kubectl get pods -l app=oauth2-server`
5. Check service status: `kubectl get svc oauth2-server`
6. Run the deployment check: `./check-deployment.sh`

Common issues:
- If the service is not accessible, check that the k3d loadbalancer is running: `docker ps | grep k3d`
- If you get a 404 error, verify that the pod is running and the service is configured correctly
- If you need to recreate the cluster, run `./manage-cluster.sh delete` followed by `./manage-cluster.sh create`
