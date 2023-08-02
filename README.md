# Go JWKS Validator

This project is a simple web server that validates JSON Web Tokens (JWT) using JSON Web Key Sets (JWKS). 
The server is written in Go and uses the [go-jose](https://github.com/go-jose/go-jose) library

## How it works

The server listens on a port specified by the `PORT` environment variable (default: 8080), 
and validates incoming JWTs using the JWKS obtained from the `JWKS_URL` 
environment variable (default: https://login.windows.net/common/discovery/keys).

The JWTs are expected to be in the Authorization header of incoming HTTP requests. 
The name of the header can be customized using the `AUTH_HEADER_NAME` environment variable 
(default: "Authorization").

The server can optionally return the validated JWT in the HTTP response.
This behavior can be controlled by the `AUTH_HEADER_RETURN` environment variable.
If `AUTH_HEADER_RETURN` is set to `false`, the server will not return the JWT. By default, it returns the JWT.

The server uses an in-memory cache to store the JWKS and reduce the number of outgoing HTTP requests. 
The JWKS is refreshed every 5 minutes.

## Usage

To run the server, you need to have Go installed. Then you can build and run the server like this:

```bash
go build
./<binary-name>
```

```bash
export PORT=8080
export JWKS_URL="https://example.com/path/to/jwks"
export AUTH_HEADER_NAME="Authorization"
export AUTH_HEADER_RETURN="true"
export SEND_BACK_CLAIMS="true"
export CACHE_TTL=300
```