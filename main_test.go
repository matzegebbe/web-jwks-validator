package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

func TestValidateTokenWithJWKS(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, edKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name      string
		algorithm jose.SignatureAlgorithm
		key       any
		claims    jwt.Claims
		want      int
	}{
		{"HMAC", jose.HS256, []byte("test-only-shared-key-at-least-32-bytes"), jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}}, http.StatusOK},
		{"RSA", jose.RS256, rsaKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}}, http.StatusOK},
		{"RSA-PSS", jose.PS256, rsaKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}}, http.StatusOK},
		{"ECDSA", jose.ES256, ecKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}}, http.StatusOK},
		{"EdDSA", jose.EdDSA, edKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}}, http.StatusOK},
		{"expired", jose.RS256, rsaKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"audience"}, Expiry: jwt.NewNumericDate(time.Now().Add(-time.Hour))}, http.StatusUnauthorized},
		{"wrong issuer", jose.RS256, rsaKey, jwt.Claims{Issuer: "other", Audience: jwt.Audience{"audience"}}, http.StatusUnauthorized},
		{"wrong audience", jose.RS256, rsaKey, jwt.Claims{Issuer: "issuer", Audience: jwt.Audience{"other"}}, http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// The handler shares a process-wide JWKS cache; keep these cases serial.
			jwksCache = keySetCache{}
			t.Cleanup(func() { jwksCache = keySetCache{} })
			key := jose.JSONWebKey{Key: tt.key, KeyID: "test", Algorithm: string(tt.algorithm), Use: "sig"}
			verificationKey := key.Public()
			if _, symmetric := tt.key.([]byte); symmetric {
				verificationKey = key
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{verificationKey}})
			}))
			t.Cleanup(server.Close)
			signer, err := jose.NewSigner(jose.SigningKey{Algorithm: tt.algorithm, Key: key}, nil)
			if err != nil {
				t.Fatal(err)
			}
			token, err := jwt.Signed(signer).Claims(tt.claims).Serialize()
			if err != nil {
				t.Fatal(err)
			}
			config := validatorConfig{jwksURL: server.URL, authHeaderName: "Authorization", expectedIssuer: "issuer", expectedAudience: "audience", cacheTTL: time.Minute, sendAccessTokenBack: true, sendAccessTokenHeaderName: "X-Auth-Access"}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			response := httptest.NewRecorder()
			validateToken(config)(response, req)
			if response.Code != tt.want {
				t.Fatalf("status = %d, want %d: %s", response.Code, tt.want, response.Body.String())
			}
			if tt.want == http.StatusOK && response.Header().Get("X-Auth-Access") != token {
				t.Fatal("validated access token was not returned")
			}
			if tt.want == http.StatusOK {
				parts := strings.Split(token, ".")
				parts[1] = base64.RawURLEncoding.EncodeToString([]byte(`{"iss":"issuer","aud":"audience","sub":"tampered"}`))
				req.Header.Set("Authorization", "Bearer "+strings.Join(parts, "."))
				response = httptest.NewRecorder()
				validateToken(config)(response, req)
				if response.Code != http.StatusUnauthorized || response.Header().Get("X-Auth-Access") != "" {
					t.Fatal("tampered token was accepted or returned")
				}
			}
		})
	}
}

func TestGetJWKSRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys":[]}`))
		_, _ = w.Write([]byte(strings.Repeat(" ", maxJWKSResponseSize)))
	}))
	t.Cleanup(server.Close)

	if _, err := getJWKS(t.Context(), server.URL); err == nil {
		t.Fatal("getJWKS() error = nil, want oversized response error")
	}
}

func TestGetJWKSUsesRequestContext(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys":[]}`))
	}))
	t.Cleanup(server.Close)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := getJWKS(ctx, server.URL)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("getJWKS() error = %v, want context.Canceled", err)
	}
}
