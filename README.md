# web-jwks-validator

![license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)

`web-jwks-validator` is a simple web server that validates
JSON Web Tokens (JWT) using JSON Web Key Sets (JWKS).
The server is written in Go and
uses the [go-jose](https://github.com/go-jose/go-jose) library.
The lib was recommended by [jwt.io](https://jwt.io)

`web-jwks-validator` is designed to be deployed within a Kubernetes
cluster, functioning as an external authentication token checker
for `ingress-nginx`. In this setup,
it validates JWTs by leveraging the
[External Oauth Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/oauth-external-auth/)
feature of `ingress-nginx`.

## Features

- JWKS fetching with cache
- JWT validation with signature verification
- Token expiration validation (exp, nbf, iat claims)
- Optional issuer and audience validation
- Custom claims checking
- Environment variables for configuration
- JWT claim details in the response (optional)

## Getting Started

1. Download or clone this repository.

   ```bash
   git clone https://github.com/matzegebbe/web-jwks-validator.git
   ```

2. Go to the project directory.

   ```bash
   cd web-jwks-validator
   ```

3. Build the project.

   ```bash
   go mod download
   go build -o web-jwks-validator
   ```

4. Start the server.

   ```bash
   export PORT=8080
   export JWKS_URL="https://login.windows.net/common/discovery/keys"
   export AUTH_HEADER_NAME="Authorization"
   export AUTH_HEADER_RETURN="true"
   export SEND_ACCESS_TOKEN_HEADER_NAME="x-auth-access"
   export SEND_BACK_CLAIMS="true"
   export CACHE_TTL=300
   export CLAIMS_CONTAINS=roles=Data.Writer
   # Optional: validate issuer and audience
   # export EXPECTED_ISSUER="https://login.microsoftonline.com/{tenant-id}/v2.0"
   # export EXPECTED_AUDIENCE="your-client-id"
   ./web-jwks-validator
   ```

## docker build

```
docker build . -t web-jwks-validator
```

## docker run

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e JWKS_URL="https://login.windows.net/common/discovery/keys" \
  -e AUTH_HEADER_NAME="Authorization" \
  -e AUTH_HEADER_RETURN="true" \
  -e SEND_ACCESS_TOKEN_HEADER_NAME="x-auth-access" \
  -e SEND_BACK_CLAIMS="true" \
  -e CACHE_TTL=300 \
  -e CLAIMS_CONTAINS="roles=Data.Writer" \
  -e EXPECTED_ISSUER="https://login.microsoftonline.com/{tenant-id}/v2.0" \
  -e EXPECTED_AUDIENCE="your-client-id" \
  ghcr.io/matzegebbe/web-jwks-validator:main
```

## Configuration

The server can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | The port the server should run on | `8080` |
| `SERVER_PATH` | The HTTP path to listen on | `/` |
| `JWKS_URL` | The URL from where the JWKS should be fetched (must be HTTPS in production) | `https://login.windows.net/common/discovery/keys` |
| `AUTH_HEADER_NAME` | The name of the header field that contains the JWT | `Authorization` |
| `AUTH_HEADER_RETURN` | If set to true, the validated access token will be included in the response headers | `true` |
| `SEND_ACCESS_TOKEN_HEADER_NAME` | If AUTH_HEADER_RETURN is on, this header will be used to send the token back (useful if downstream needs the token) | `Authorization` |
| `SEND_BACK_CLAIMS` | If set to true, all claims from the JWT will be returned as JSON in the response | `true` |
| `CACHE_TTL` | Defines how long JWKS should be cached in seconds | `300` |
| `CLAIMS_CONTAINS` | A comma-separated list of required claim key=value pairs | *(not checked)* |
| `EXPECTED_ISSUER` | If set, validates that the token's `iss` claim matches this value | *(not checked)* |
| `EXPECTED_AUDIENCE` | If set, validates that the token's `aud` claim matches this value | *(not checked)* |

### Security Validation

The validator performs the following checks on each token:

1. **Signature verification** - Token must be signed by a key from the JWKS
2. **Expiration check** - Token must not be expired (`exp` claim)
3. **Not before check** - Token must be valid for use (`nbf` claim)
4. **Issuer validation** - If `EXPECTED_ISSUER` is set, the `iss` claim must match
5. **Audience validation** - If `EXPECTED_AUDIENCE` is set, the `aud` claim must match
6. **Custom claims** - If `CLAIMS_CONTAINS` is set, all specified claims must be present with matching values

## ingress-nginx check

### ingress-nginx with auth annotations
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  annotations:
    nginx.ingress.kubernetes.io/enable-global-auth: false
    nginx.ingress.kubernetes.io/auth-url: http://web-jwks-validator-service.services.svc.cluster.local
    nginx.ingress.kubernetes.io/auth-response-headers: Authorization
    nginx.ingress.kubernetes.io/auth-snippet: |
      access_log off;
spec:
  ingressClassName: public
  rules:
    - host: service.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: backend-service
                port:
                  number: 80
  tls:
    - hosts: [service.example.com]
      secretName: 'ssl-cert'
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-jwks-validator-deployment
  labels:
    app: web-jwks-validator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: web-jwks-validator
  template:
    metadata:
      labels:
        app: web-jwks-validator
    spec:
      containers:
        - name: web-jwks-validator
          image: ghcr.io/matzegebbe/web-jwks-validator:main-5767335050
          ports:
            - containerPort: 8080
          env:
            - name: JWKS_URL
              value: https://login.microsoftonline.com/common/discovery/v2.0/keys
            # Optional: validate issuer and audience for additional security
            # - name: EXPECTED_ISSUER
            #   value: https://login.microsoftonline.com/{tenant-id}/v2.0
            # - name: EXPECTED_AUDIENCE
            #   value: your-client-id
          resources:
           requests:
            cpu: "100m"
            memory: "128Mi"
           limits:
            cpu: "200m"
            memory: "128Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: web-jwks-validator-service
spec:
  selector:
    app: web-jwks-validator
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

## Test

### Quick Test with Duende Demo (No Setup Required)

The easiest way to test the validator is using the public [Duende IdentityServer demo](https://demo.duendesoftware.com). No credentials or setup required.

**1. Start the validator with Docker:**

```bash
docker run --rm -p 8080:8080 \
  -e JWKS_URL="https://demo.duendesoftware.com/.well-known/openid-configuration/jwks" \
  ghcr.io/matzegebbe/web-jwks-validator:main
```

**2. Run the test script:**

```bash
./misc/test_duende_demo.sh
```

Or test manually with curl:

```bash
# Get a token from the Duende demo
TOKEN=$(curl -s -X POST https://demo.duendesoftware.com/connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=m2m&client_secret=secret&grant_type=client_credentials&scope=api" \
  | jq -r '.access_token')

# Test against the validator
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/
```

**All-in-one Docker test command:**

```bash
./misc/all_in_one_demo.sh
```

This script starts the validator in Docker, fetches a token from Duende, tests the validation, and cleans up automatically.

### Test with Microsoft App Registration

For testing against Azure AD, use the [test_call bash script](misc/test_call.sh):

```bash
./misc/test_call.sh <CLIENT_ID> <CLIENT_SECRET> <TENANT_ID>
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [go-jose](https://github.com/go-jose/go-jose) library

Enjoy using `web-jwks-validator`! Your feedback and contributions are welcome.
