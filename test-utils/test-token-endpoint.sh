#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Testing /token endpoint..."

# Test 1: Missing Authorization header
echo -e "\n${GREEN}Test 1: Missing Authorization header${NC}"
curl -X POST http://localhost:8080/token

# Test 2: Invalid Basic Auth format
echo -e "\n\n${GREEN}Test 2: Invalid Basic Auth format${NC}"
curl -X POST http://localhost:8080/token -H "Authorization: Bearer token"

# Test 3: Invalid base64 encoding
echo -e "\n\n${GREEN}Test 3: Invalid base64 encoding${NC}"
curl -X POST http://localhost:8080/token -H "Authorization: Basic invalid-base64"

# Test 4: Valid request
echo -e "\n\n${GREEN}Test 4: Valid request${NC}"
curl -X POST http://localhost:8080/token \
  -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
  -H "Content-Type: application/json"

echo -e "\n\nTests completed!"
