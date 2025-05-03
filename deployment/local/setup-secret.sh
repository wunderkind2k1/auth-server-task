#!/bin/bash

set -e

# Get the first private key from the keytool directory
PRIVATE_KEY_PATH=$(ls -1 ../../keytool/keys/*.private.pem | head -n 1)
if [ -z "$PRIVATE_KEY_PATH" ]; then
    echo "No private key found in keytool/keys directory"
    exit 1
fi

# Read the private key content
PRIVATE_KEY_CONTENT=$(cat "$PRIVATE_KEY_PATH")
if [ -z "$PRIVATE_KEY_CONTENT" ]; then
    echo "Failed to read private key content"
    exit 1
fi

# Create the Kubernetes secret
echo "Creating Kubernetes secret from key: $PRIVATE_KEY_PATH"
kubectl create secret generic jwt-key \
    --from-literal=private-key="$PRIVATE_KEY_CONTENT" \
    --dry-run=client -o yaml | kubectl apply -f -

echo "Secret created successfully"
