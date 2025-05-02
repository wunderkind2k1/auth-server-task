#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Testing /introspect endpoint..."

# Test 1: Missing token
echo -e "\n${GREEN}Test 1: Missing token${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/introspect)
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "400" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

# Test 2: Invalid HTTP method
echo -e "\n\n${GREEN}Test 2: Invalid HTTP method${NC}"
response=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/introspect)
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "405" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body"
fi

# Test 3: Invalid token
echo -e "\n\n${GREEN}Test 3: Invalid token${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/introspect \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=invalid.token.here")
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "200" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
    if ! echo "$body" | jq -e '.active == false' > /dev/null; then
        echo -e "${RED}Expected active: false for invalid token${NC}"
    fi
fi

# Test 4: Valid token (using token from /token endpoint)
echo -e "\n\n${GREEN}Test 4: Valid token${NC}"
# First get a valid token
token_response=$(curl -s -X POST http://localhost:8080/token \
  -H "Authorization: Basic $(echo -n 'sho:test123' | base64)" \
  -H "Content-Type: application/json")
access_token=$(echo "$token_response" | jq -r '.access_token')

# Now test introspection with the valid token
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/introspect \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "token=$access_token")
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "200" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'

    # Verify required fields
    required_fields=("active" "token_type" "sub" "iss" "exp" "iat" "nbf")
    for field in "${required_fields[@]}"; do
        if ! echo "$body" | jq -e ".$field" > /dev/null; then
            echo -e "${RED}Response is missing required field: $field${NC}"
        else
            echo -e "${GREEN}✓${NC} Field '$field' is present"
        fi
    done
fi

echo -e "\n\nTests completed!"
