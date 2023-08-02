#!/bin/bash

CLIENT_ID=$1
CLIENT_SECRET=$2
TENANT_ID=$3

RESULT=$(curl --silent -X POST https://login.microsoftonline.com/${TENANT_ID}/oauth2/v2.0/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "client_id=${CLIENT_ID}" \
  --data-urlencode "scope=${CLIENT_ID}/.default" \
  --data-urlencode "client_secret=${CLIENT_SECRET}" \
  --data-urlencode "grant_type=client_credentials")

echo "${RESULT}"

ACCESS_TOKEN=$(echo "${RESULT}" | jq .access_token)

jwt decode "${ACCESS_TOKEN//\"}"

echo "${ACCESS_TOKEN//\"}"
# JSON Web Key Sets endpoint https://sts.windows.net/${TENANT_ID}/.well-known/openid-configuration

curl -v -H "Authorization: Bearer "${ACCESS_TOKEN//\"} http://localhost:8080
