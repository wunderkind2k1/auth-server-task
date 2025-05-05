#!/usr/bin/env bash

set -euo pipefail

# Check if REPLICAS is set
if [ -z "${REPLICAS:-}" ]; then
    echo "Error: REPLICAS environment variable is not set"
    echo "Usage: REPLICAS=2 ./scale.sh"
    exit 1
fi

echo "Scaling OAuth2 server to ${REPLICAS} replicas..."

# Scale the deployment
kubectl scale deployment oauth2-server --replicas="${REPLICAS}" -n default

echo "OAuth2 server scaled successfully to ${REPLICAS} replicas."
