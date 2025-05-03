#!/bin/bash

set -e

IMAGE_NAME="oauth2-server"
IMAGE_TAG="latest"
TAR_FILE="oauth2-server.tar"

echo "Building Docker image: $IMAGE_NAME:$IMAGE_TAG"
# Change to project root for the build
cd ../..
docker build -t $IMAGE_NAME:$IMAGE_TAG \
    -f deployment/local/Dockerfile .

echo "Saving Docker image to $TAR_FILE"
docker save $IMAGE_NAME:$IMAGE_TAG > deployment/local/$TAR_FILE

echo "Image built and saved successfully"
