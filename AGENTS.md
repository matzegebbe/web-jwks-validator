# AGENTS.md

## Project overview
`web-jwks-validator` is a Go web server that validates JWTs against a JWKS endpoint. It is designed to run as an auth service for ingress-nginx and exposes HTTP endpoints for token validation.

## Repository layout
- `main.go`: application entry point and HTTP server wiring.
- `config.go`: environment configuration parsing.
- `claim_test.go`: JWT claim validation tests.
- `misc/`: helper scripts (e.g., `misc/test_call.sh`).
- `Dockerfile`, `Dockerfile-goreleaser`: container build definitions.
- `release-please-config.json`, `CHANGELOG.md`, `RELEASE.md`: release automation and changelog.

## Build and run
### Local build
```bash
go mod download
go build -o web-jwks-validator
```

### Local run
```bash
export PORT=8080
export JWKS_URL="https://login.windows.net/common/discovery/keys"
export AUTH_HEADER_NAME="Authorization"
export AUTH_HEADER_RETURN="true"
export SEND_ACCESS_TOKEN_HEADER_NAME="x-auth-access"
export SEND_BACK_CLAIMS="true"
export CACHE_TTL=300
export CLAIMS_CONTAINS="roles=Data.Writer"
# Optional: validate issuer and audience
# export EXPECTED_ISSUER="https://login.microsoftonline.com/{tenant-id}/v2.0"
# export EXPECTED_AUDIENCE="your-client-id"
./web-jwks-validator
```

### Docker build/run
```bash
docker build . -t web-jwks-validator
```

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
  ghcr.io/matzegebbe/web-jwks-validator:main
```

## Conventions
- Use **Conventional Commits** for all commit messages.
- Follow **Semantic Versioning (SemVer)** for releases.
