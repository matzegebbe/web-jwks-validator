package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

// Cache struct
type Cache struct {
	sync.RWMutex
	jwks    *jose.JSONWebKeySet
	expires time.Time
}

type CustomClaims struct {
	Roles []string `json:"roles,omitempty"`
}

var (
	jwksCache = Cache{}
)

func main() {

	port := GetPortFromEnv()
	jwksUrl := GetJwksURLFromEnv()
	authHeaderName := GetAuthHeaderNameFromEnv()
	sendAccessTokenBack := GetAccessTokenSettingFromEnv()
	sendAllClaimsAsJson := GetSendAllClaimsAsJson()
	ttlInSeconds := GetTTLFromEnv()
	claimContainsCheck := GetClaimContains()

	http.HandleFunc(GetPathFromEnv(), validateToken(
		jwksUrl,
		authHeaderName,
		sendAccessTokenBack,
		ttlInSeconds,
		sendAllClaimsAsJson,
		claimContainsCheck,
	))

	http.ListenAndServe(":"+port, nil)
}
func validateToken(
	jwksURL string,
	authHeaderName string,
	sendAccessTokenBack bool,
	ttlInSeconds int,
	sendAllClaimsAsJson bool,
	claimContainsCheck []string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r, authHeaderName)
		if tokenString == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseSigned(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		keys, err := getJwksWithCache(jwksURL, ttlInSeconds)
		if err != nil {
			http.Error(w, "Error fetching JWKS", http.StatusInternalServerError)
			return
		}

		claims := make(map[string]interface{})
		for _, key := range keys.Keys {
			err = token.Claims(key, &claims)
			if err == nil {
				if len(claimContainsCheck) > 0 {
					if !checkIfClaimContainsAllClaimContainsCheck(claims, claimContainsCheck) {
						http.Error(w, "Missing required claims from token", http.StatusUnauthorized)
					}
				}
				if sendAccessTokenBack {
					w.Header().Set(authHeaderName, tokenString)
				}
				if sendAllClaimsAsJson {
					responseJSON, err := json.Marshal(claims)
					if err != nil {
						http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
						return
					}
					_, err = w.Write(responseJSON)
					if err != nil {
						http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
						return
					}
				} else {
					w.Write([]byte("Token valid."))
				}
				return
			}
		}

		http.Error(w, "Token could not be validated", http.StatusUnauthorized)
	}

}

func checkIfClaimContainsAllClaimContainsCheck(claims map[string]interface{}, claimContainsCheck []string) bool {
	for _, claimCheck := range claimContainsCheck {
		claimCheck = strings.ReplaceAll(claimCheck, "\"", "")
		parts := strings.Split(claimCheck, "=")
		if len(parts) != 2 {
			log.Println("Invalid claim check", claimCheck)
			continue
		}
		key, value := parts[0], parts[1]
		claimValue, exists := claims[key]
		if !exists {
			continue
		}

		switch v := claimValue.(type) {
		case string:
			if v != value {
				// the string value does not match
				return false
			}
		case []string:
			found := false
			for _, s := range v {
				if s == value {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case []interface{}:
			found := false
			for _, iface := range v {
				s, ok := iface.(string)
				if !ok {
					log.Printf("Unexpected type %T in []interface{} for key %s\n", iface, key)
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
	authHeader := r.Header.Get(authHeaderName)
	// often we have bearer token
	strArr := strings.Split(authHeader, " ")
	if len(strArr) == 2 {
		return strArr[1]
	} else if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
}

func getJwksWithCache(jwksURL string, ttlInSeconds int) (*jose.JSONWebKeySet, error) {
	jwksCache.RLock()
	if jwksCache.jwks != nil && time.Now().Before(jwksCache.expires) {
		defer jwksCache.RUnlock()
		return jwksCache.jwks, nil
	}
	jwksCache.RUnlock()

	jwksCache.Lock()
	defer jwksCache.Unlock()

	// Double-checked locking
	if jwksCache.jwks != nil && time.Now().Before(jwksCache.expires) {
		return jwksCache.jwks, nil
	}

	jwks, err := getJwks(jwksURL)
	if err != nil {
		return nil, err
	}

	jwksCache.jwks = jwks
	jwksCache.expires = time.Now().Add(time.Duration(ttlInSeconds) * time.Second)

	return jwksCache.jwks, nil
}

func getJwks(jwksURL string) (*jose.JSONWebKeySet, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks = new(jose.JSONWebKeySet)
	err = json.NewDecoder(resp.Body).Decode(jwks)

	return jwks, err
}
