#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Testing /.well-known/jwks.json endpoint..."

# Test 1: GET request
echo -e "\n${GREEN}Test 1: GET request${NC}"
response=$(curl -s -w "\n%{http_code}" -X GET http://localhost:8080/.well-known/jwks.json)
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "200" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
else
    echo "Response:"
    echo "$body" | jq '.'
fi

# Test 2: POST request (should fail)
echo -e "\n\n${GREEN}Test 2: POST request (should fail)${NC}"
response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/.well-known/jwks.json)
status_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')
if [ "$status_code" != "405" ]; then
    echo -e "${RED}Unexpected status code: $status_code${NC}"
    echo -e "${RED}Response: $body${NC}"
    exit 1
fi

if [ "$body" != "Method not allowed: POST" ]; then
    echo -e "${RED}Unexpected error message: $body${NC}"
    echo -e "${RED}Expected: Method not allowed: POST${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Status code: $status_code"
echo -e "${GREEN}✓${NC} Error message: $body"

# Test 3: Verify JWKS structure
echo -e "\n\n${GREEN}Test 3: Verify JWKS structure${NC}"
response=$(curl -s http://localhost:8080/.well-known/jwks.json)
if ! echo "$response" | jq -e '.keys' > /dev/null; then
    echo -e "${RED}Response does not contain 'keys' array${NC}"
    echo -e "${RED}Response: $response${NC}"
else
    echo "JWKS structure is valid:"
    echo "$response" | jq '.'

    # Verify each key has required fields
    keys_count=$(echo "$response" | jq '.keys | length')
    echo -e "\nFound $keys_count key(s)"

    for i in $(seq 0 $(($keys_count-1))); do
        key=$(echo "$response" | jq ".keys[$i]")
        echo -e "\nVerifying key $((i+1)):"

        # Check required fields
        required_fields=("kty" "use" "kid" "n" "e")
        for field in "${required_fields[@]}"; do
            if ! echo "$key" | jq -e ".$field" > /dev/null; then
                echo -e "${RED}Key $((i+1)) is missing required field: $field${NC}"
            else
                echo -e "${GREEN}✓${NC} Field '$field' is present"
            fi
        done
    done
fi

echo -e "\n\nTests completed!"
