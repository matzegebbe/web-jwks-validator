# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`web-jwks-validator` is a Go web server that validates JWTs against a JWKS endpoint. It runs as an external authentication service for ingress-nginx in Kubernetes clusters, using the [go-jose](https://github.com/go-jose/go-jose) library.

## Build and Test Commands

```bash
# Build
go mod download
go build -o web-jwks-validator

# Run tests
go test ./...

# Docker build
docker build . -t web-jwks-validator
```

## Architecture

The codebase is minimal (~375 lines of Go) with a single external dependency (go-jose).

| File | Purpose |
|------|---------|
| `main.go` | HTTP server, JWT validation logic, JWKS caching with double-checked locking |
| `config.go` | Environment variable parsing |
| `claim_test.go` | Unit tests for claim validation |

**HTTP Handler Flow:**
1. `extractToken()` - Extract JWT from Authorization header
2. `jwt.ParseSigned()` - Parse the token
3. `getJwksWithCache()` - Fetch JWKS with TTL-based caching
4. Validate signature against JWKS keys
5. Validate time-based claims (exp, nbf, iat) and optionally issuer/audience
6. `checkIfClaimContainsAllClaimContainsCheck()` - Validate required claims
7. Return claims as JSON or success message

**Key Configuration (via environment variables):**
- `PORT` (default: 8080)
- `JWKS_URL` - JWKS endpoint URL (must be HTTPS)
- `AUTH_HEADER_NAME` (default: Authorization)
- `CACHE_TTL` - JWKS cache TTL in seconds (default: 300)
- `CLAIMS_CONTAINS` - Required claims as `key=value,key=value`
- `EXPECTED_ISSUER` - Validate token issuer (optional)
- `EXPECTED_AUDIENCE` - Validate token audience (optional)

## Conventions

- Use **Conventional Commits** for all commit messages
- Follow **Semantic Versioning (SemVer)** for releases
- Releases are automated via Release Please
