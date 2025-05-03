#!/bin/bash

set -e

echo "Setting up port-forwarding for oauth2-server..."
echo "The service will be available at http://localhost:8080"
echo "Press Ctrl+C to stop port-forwarding"
echo ""

kubectl port-forward svc/oauth2-server 8080:80
