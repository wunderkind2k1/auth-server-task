#!/bin/bash

set -e

CLUSTER_NAME="oauth2-cluster"
CLUSTER_PORT=8080

function create_cluster() {
    echo "Creating k3d cluster: $CLUSTER_NAME"
    k3d cluster create $CLUSTER_NAME \
        --port "$CLUSTER_PORT:8080@loadbalancer" \
        --wait
    echo "Cluster created successfully"
}

function delete_cluster() {
    echo "Deleting k3d cluster: $CLUSTER_NAME"
    k3d cluster delete $CLUSTER_NAME
    echo "Cluster deleted successfully"
}

function check_cluster() {
    if k3d cluster list | grep -q "$CLUSTER_NAME"; then
        return 0
    else
        return 1
    fi
}

case "$1" in
    "create")
        if check_cluster; then
            echo "Cluster already exists"
            exit 0
        fi
        create_cluster
        ;;
    "delete")
        if ! check_cluster; then
            echo "Cluster does not exist"
            exit 0
        fi
        delete_cluster
        ;;
    "status")
        if check_cluster; then
            echo "Cluster is running"
        else
            echo "Cluster is not running"
        fi
        ;;
    *)
        echo "Usage: $0 {create|delete|status}"
        exit 1
        ;;
esac
