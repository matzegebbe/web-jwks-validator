#!/bin/bash
#
# Test script using the public Duende IdentityServer demo
# No credentials or setup required - great for quick testing
#
# Usage:
#   ./test_duende_demo.sh [validator_url]
#
# Examples:
#   ./test_duende_demo.sh                      # Test against localhost:8080
#   ./test_duende_demo.sh http://localhost:9000  # Test against custom URL
#

set -e

VALIDATOR_URL="${1:-http://localhost:8080}"
DUENDE_TOKEN_URL="https://demo.duendesoftware.com/connect/token"
DUENDE_JWKS_URL="https://demo.duendesoftware.com/.well-known/openid-configuration/jwks"

echo "=== Duende Demo Test ==="
echo "JWKS URL: ${DUENDE_JWKS_URL}"
echo "Validator URL: ${VALIDATOR_URL}"
echo ""

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    echo "Install with: brew install jq (macOS) or apt-get install jq (Linux)"
    exit 1
fi

# Get token from Duende demo
echo "Fetching token from Duende demo..."
RESULT=$(curl --silent -X POST "${DUENDE_TOKEN_URL}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=m2m" \
  -d "client_secret=secret" \
  -d "grant_type=client_credentials" \
  -d "scope=api")

ACCESS_TOKEN=$(echo "${RESULT}" | jq -r '.access_token')

if [ "${ACCESS_TOKEN}" == "null" ] || [ -z "${ACCESS_TOKEN}" ]; then
    echo "Error: Failed to get access token"
    echo "Response: ${RESULT}"
    exit 1
fi

echo "Token received successfully"
echo ""

# Decode token header and payload (without verification)
echo "=== Token Info ==="
echo "${ACCESS_TOKEN}" | cut -d'.' -f2 | base64 -d 2>/dev/null | jq . 2>/dev/null || echo "(could not decode payload)"
echo ""

# Test against validator
echo "=== Testing Validator ==="
echo "Sending request to ${VALIDATOR_URL}..."
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  "${VALIDATOR_URL}")

echo "HTTP Status: ${HTTP_CODE}"
echo "Response:"
cat /tmp/response.txt
echo ""

if [ "${HTTP_CODE}" == "200" ]; then
    echo ""
    echo "=== SUCCESS ==="
else
    echo ""
    echo "=== FAILED ==="
    exit 1
fi
