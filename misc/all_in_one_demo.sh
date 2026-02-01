#!/bin/bash
#
# All-in-one demo script for web-jwks-validator
# Starts the validator in Docker, tests with a real JWT, and cleans up
#
# No setup required - uses the public Duende IdentityServer demo
#
# Usage:
#   ./all_in_one_demo.sh
#

set -e

CONTAINER_NAME="jwks-validator-demo"
DUENDE_JWKS_URL="https://demo.duendesoftware.com/.well-known/openid-configuration/jwks"
DUENDE_TOKEN_URL="https://demo.duendesoftware.com/connect/token"
IMAGE="ghcr.io/matzegebbe/web-jwks-validator:main"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed."
    echo "Install with: brew install jq (macOS) or apt-get install jq (Linux)"
    exit 1
fi

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up..."
    docker stop "${CONTAINER_NAME}" 2>/dev/null || true
}
trap cleanup EXIT

echo "=== web-jwks-validator All-in-One Demo ==="
echo ""

# Start validator
echo "Starting validator container..."
podman run -d --rm --name "${CONTAINER_NAME}" -p 8080:8080 \
  -e JWKS_URL="${DUENDE_JWKS_URL}" \
  "${IMAGE}"

echo "Waiting for container to be ready..."
sleep 2

# Get token
echo "Fetching token from Duende demo..."
TOKEN=$(curl -s -X POST "${DUENDE_TOKEN_URL}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=m2m&client_secret=secret&grant_type=client_credentials&scope=api" \
  | jq -r '.access_token')

if [ "${TOKEN}" == "null" ] || [ -z "${TOKEN}" ]; then
    echo "Error: Failed to get access token"
    exit 1
fi

echo "Token received successfully"
echo ""

# Test validator
echo "=== Testing Validator ==="
echo ""

HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" \
  -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:8080/)

echo "HTTP Status: ${HTTP_CODE}"
echo "Response:"
cat /tmp/response.txt | jq . 2>/dev/null || cat /tmp/response.txt
echo ""

if [ "${HTTP_CODE}" == "200" ]; then
    echo "=== SUCCESS ==="
    echo "The validator correctly validated the JWT against the JWKS endpoint"
else
    echo "=== FAILED ==="
    exit 1
fi
