package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type keySetCache struct {
	mu        sync.RWMutex
	keySet    *jose.JSONWebKeySet
	expiresAt time.Time
}

type validatorConfig struct {
	jwksURL                   string
	authHeaderName            string
	sendAccessTokenBack       bool
	sendAccessTokenHeaderName string
	cacheTTL                  time.Duration
	sendAllClaimsAsJSON       bool
	requiredClaims            []string
	expectedIssuer            string
	expectedAudience          string
}

var (
	jwksCache = keySetCache{}
	version   = "dev"
	commit    = "none"
	date      = "unknown"
)

func main() {
	log.Printf("VERSION %s, COMMIT %s, BUILD AT %s", version, commit, date)

	config := validatorConfig{
		jwksURL:                   getJWKSURLFromEnv(),
		authHeaderName:            getAuthHeaderNameFromEnv(),
		sendAccessTokenBack:       getSendBackAccessTokenFromEnv(),
		sendAccessTokenHeaderName: getSendBackAccessTokenNameFromEnv(),
		cacheTTL:                  time.Duration(getCacheTTLFromEnv()) * time.Second,
		sendAllClaimsAsJSON:       getSendAllClaimsAsJSONFromEnv(),
		requiredClaims:            getRequiredClaimsFromEnv(),
		expectedIssuer:            getExpectedIssuerFromEnv(),
		expectedAudience:          getExpectedAudienceFromEnv(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(getPathFromEnv(), validateToken(config))

	server := &http.Server{
		Addr:              ":" + getPortFromEnv(),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func validateToken(config validatorConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r, config.authHeaderName)
		if tokenString == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseSigned(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		keys, err := getJWKSWithCache(r.Context(), config.jwksURL, config.cacheTTL)
		if err != nil {
			http.Error(w, "Error fetching JWKS", http.StatusInternalServerError)
			return
		}

		for _, key := range keys.Keys {
			claims := make(map[string]any)
			var standardClaims jwt.Claims
			if err := token.Claims(key, &claims, &standardClaims); err != nil {
				continue
			}

			expected := jwt.Expected{Time: time.Now()}
			if config.expectedIssuer != "" {
				expected.Issuer = config.expectedIssuer
			}
			if config.expectedAudience != "" {
				expected.Audience = jwt.Audience{config.expectedAudience}
			}

			if err := standardClaims.Validate(expected); err != nil {
				http.Error(w, "Token validation failed", http.StatusUnauthorized)
				return
			}
			if len(config.requiredClaims) > 0 && !claimsContainAll(claims, config.requiredClaims) {
				http.Error(w, "Missing required claims from token", http.StatusUnauthorized)
				return
			}
			if config.sendAccessTokenBack {
				w.Header().Set(config.sendAccessTokenHeaderName, tokenString)
			}
			if config.sendAllClaimsAsJSON {
				responseJSON, err := json.Marshal(claims)
				if err != nil {
					http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
					return
				}
				if _, err := w.Write(responseJSON); err != nil {
					http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
					return
				}
			} else {
				_, _ = w.Write([]byte("Token valid."))
			}
			return
		}

		http.Error(w, "Token could not be validated", http.StatusUnauthorized)
	}

}

func claimsContainAll(claims map[string]any, requiredClaimChecks []string) bool {
	for _, claimCheck := range requiredClaimChecks {
		claimCheck = strings.ReplaceAll(claimCheck, "\"", "")
		key, value, ok := strings.Cut(claimCheck, "=")
		if !ok {
			log.Println("Invalid claim check", claimCheck)
			return false
		}
		claimValue, exists := claims[key]
		if !exists {
			return false
		}

		switch v := claimValue.(type) {
		case string:
			if v != value {
				// the string value does not match
				return false
			}
		case []string:
			if !slices.Contains(v, value) {
				return false
			}
		case []any:
			found := false
			for _, item := range v {
				s, ok := item.(string)
				if !ok {
					log.Printf("Unexpected type %T in []any for key %s\n", item, key)
					return false
				}
				if s == value {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		default:
			log.Printf("Unexpected type %T for key %s\n", v, key)
			return false
		}
	}
	return true
}

func extractToken(r *http.Request, authHeaderName string) string {
	fields := strings.Fields(r.Header.Get(authHeaderName))
	switch len(fields) {
	case 1:
		return fields[0]
	case 2:
		return fields[1]
	default:
		return ""
	}
}

func getJWKSWithCache(ctx context.Context, jwksURL string, ttl time.Duration) (*jose.JSONWebKeySet, error) {
	jwksCache.mu.RLock()
	if jwksCache.keySet != nil && time.Now().Before(jwksCache.expiresAt) {
		defer jwksCache.mu.RUnlock()
		return jwksCache.keySet, nil
	}
	jwksCache.mu.RUnlock()

	jwksCache.mu.Lock()
	defer jwksCache.mu.Unlock()

	// Double-checked locking
	if jwksCache.keySet != nil && time.Now().Before(jwksCache.expiresAt) {
		return jwksCache.keySet, nil
	}

	jwks, err := getJWKS(ctx, jwksURL)
	if err != nil {
		return nil, err
	}

	jwksCache.keySet = jwks
	jwksCache.expiresAt = time.Now().Add(ttl)

	return jwksCache.keySet, nil
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

const maxJWKSResponseSize = 1 << 20 // 1 MiB

func getJWKS(ctx context.Context, jwksURL string) (*jose.JSONWebKeySet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create JWKS request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch JWKS: unexpected HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxJWKSResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("read JWKS response: %w", err)
	}
	if len(body) > maxJWKSResponseSize {
		return nil, fmt.Errorf("read JWKS response: exceeds %d-byte limit", maxJWKSResponseSize)
	}

	jwks := new(jose.JSONWebKeySet)
	if err := json.Unmarshal(body, jwks); err != nil {
		return nil, fmt.Errorf("decode JWKS response: %w", err)
	}

	return jwks, nil
}
