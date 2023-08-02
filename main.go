package main

import (
	"encoding/json"
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

	http.HandleFunc(GetPathFromEnv(), validateToken(
		jwksUrl,
		authHeaderName,
		sendAccessTokenBack,
		ttlInSeconds,
		sendAllClaimsAsJson))

	http.ListenAndServe(":"+port, nil)
}
func validateToken(
	jwksURL string,
	authHeaderName string,
	sendAccessTokenBack bool,
	ttlInSeconds int,
	sendAllClaimsAsJson bool) http.HandlerFunc {
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

		cl := jwt.Claims{}
		for _, key := range keys.Keys {
			err = token.Claims(key, &cl)
			if err == nil {
				allClaims := make(map[string]interface{})
				if err := token.UnsafeClaimsWithoutVerification(&allClaims); err != nil {
					http.Error(w, "Failed to parse custom claims", http.StatusInternalServerError)
					return
				}
				responseJSON, err := json.Marshal(allClaims)
				if err != nil {
					http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
					return
				}
				if sendAccessTokenBack {
					w.Header().Set(authHeaderName, tokenString)
				}
				w.Write(responseJSON)
				return
			}
		}

		http.Error(w, "Token could not be validated", http.StatusUnauthorized)
	}

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
