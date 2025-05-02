#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Testing /token endpoint..."

# Test 1: Missing Authorization header
echo -e "\n${GREEN}Test 1: Missing Authorization header${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/token)
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "401" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

# Test 2: Invalid Basic Auth format
echo -e "\n\n${GREEN}Test 2: Invalid Basic Auth format${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/token -H "Authorization: Bearer token")
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "401" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

# Test 3: Invalid base64 encoding
echo -e "\n\n${GREEN}Test 3: Invalid base64 encoding${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/token -H "Authorization: Basic invalid-base64")
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "401" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

# Test 4: Valid request
echo -e "\n\n${GREEN}Test 4: Valid request${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/token \
  -H "Authorization: Basic $(echo -n 'sho:test123' | base64)" \
  -H "Content-Type: application/json")
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "200" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

echo -e "\n\nTests completed!"
