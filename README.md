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
- JWT validation with claims checking
- Environment variables for configuration
- JWT claim details in the response (optional)

## Getting Started

1. Download or clone this repository.

   ```bash
   git clone https://github.com/your-github-username/web-jwks-validator.git
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
   export SEND_BACK_CLAIMS="true"
   export CACHE_TTL=300
   export CLAIMS_CONTAINS=roles=Data.Writer
   ./web-jwks-validator
   ```

## docker build

```
docker build . -t web-jwks-validator
```

## Configuration

The server can be configured using the following environment variables:

- `PORT`: the port the server should run on
  (default *8080*)
- `JWKS_URL`: the URL from where the JWKS should be fetched
  (default *https://login.windows.net/common/discovery/keys*)
- `AUTH_HEADER_NAME`: the name of the header field that contains the JWT
  (default *Authorization*)
- `SEND_ACCESS_TOKEN_BACK`: if set to true, the validated access token will be included in the response headers
  (default *true*)
- `SEND_ACCESS_TOKEN_HEADER_NAME`: if SEND_ACCESS_TOKEN_BACK is on this header will be used to send the token back
  this can be useful if the downstream application needs to token as well
  (default *Authorization*)
- `SEND_ALL_CLAIMS_AS_JSON`: if set to true, all claims from the JWT will be returned as JSON in the response
  (default *true*)
- `TTL_IN_SECONDS`: defines how long JWKS should be cached
  (default *300*)
- `CLAIM_CONTAINS`: a comma-separated list of required claim key=value pairs
  (default *""* - not checked )

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

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [go-jose](https://github.com/go-jose/go-jose) library

Enjoy using `web-jwks-validator`! Your feedback and contributions are welcome.
