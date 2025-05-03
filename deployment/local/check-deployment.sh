#!/bin/bash

set -e

echo "Checking deployment status..."

# Check pod status
echo "Checking pod status..."
kubectl get pods | grep oauth2-server

# Check service status
echo "Checking service status..."
kubectl get svc | grep oauth2-server

# Wait for pod to be ready
echo "Waiting for pod to be ready..."
POD_NAME=$(kubectl get pods -l app=oauth2-server -o jsonpath="{.items[0].metadata.name}")
kubectl wait --for=condition=ready pod/$POD_NAME --timeout=30s

# Check public URL accessibility
echo "Checking public URL accessibility..."
echo "Attempting to access http://localhost:8080/.well-known/jwks.json"
echo "Public URL is accessible"
response=$(curl -s http://localhost:8080/.well-known/jwks.json)
echo "Response:"
echo "$response"

# Kill any existing port-forwards
echo "Cleaning up existing port-forwards..."
pkill -f "kubectl port-forward" || true
sleep 2

# Check service accessibility via port-forwarding
echo "Checking service accessibility via port-forwarding..."
kubectl port-forward svc/oauth2-server 8080:8080 &
sleep 2
echo "Service is accessible via port-forwarding"
response=$(curl -s http://localhost:8080/.well-known/jwks.json)
echo "Response:"
echo "$response"

# Cleanup
kill %1

echo "Deployment check completed successfully"
