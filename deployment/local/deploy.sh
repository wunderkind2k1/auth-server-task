#!/bin/bash

set -e

CLUSTER_NAME="oauth2-cluster"
TAR_FILE="oauth2-server.tar"
IMAGE_NAME="oauth2-server:latest"

# Check if cluster exists
if ! k3d cluster list | grep -q "$CLUSTER_NAME"; then
    echo "Cluster does not exist. Please create it first using ./manage-cluster.sh create"
    exit 1
fi

# Check if image exists in the cluster
if ! k3d image list | grep -q "$IMAGE_NAME"; then
    echo "Loading image into cluster..."
    k3d image import $TAR_FILE -c "$CLUSTER_NAME"
else
    echo "Image already exists in cluster, skipping import"
fi

# Apply the Kubernetes manifests
echo "Deploying application..."
kubectl apply -f k8s/

echo "Deployment completed successfully"
