#!/usr/bin/env bash

set -euo pipefail

echo "Undeploying OAuth2 server..."

# Delete the deployment
kubectl delete deployment oauth2-server -n default

echo "OAuth2 server undeployed successfully."
