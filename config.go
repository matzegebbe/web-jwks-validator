package main

import (
	"cmp"
	"log"
	"os"
	"strconv"
	"strings"
)

func getPathFromEnv() string {
	path := cmp.Or(os.Getenv("SERVER_PATH"), "/")
	log.Println("SERVER_PATH:", path)
	return path
}

func getPortFromEnv() string {
	port := cmp.Or(os.Getenv("PORT"), "8080")
	log.Println("PORT:", port)
	return port
}

func getJWKSURLFromEnv() string {
	jwksURL := cmp.Or(os.Getenv("JWKS_URL"), "https://login.windows.net/common/discovery/keys")
	// Warn if JWKS URL is not using HTTPS (insecure)
	if !strings.HasPrefix(jwksURL, "https://") {
		log.Println("WARNING: JWKS_URL is not using HTTPS - this is insecure and vulnerable to MITM attacks")
	}
	log.Println("JWKS_URL:", jwksURL)
	return jwksURL
}

func getAuthHeaderNameFromEnv() string {
	authHeaderName := cmp.Or(os.Getenv("AUTH_HEADER_NAME"), "Authorization")
	log.Println("AUTH_HEADER_NAME:", authHeaderName)
	return authHeaderName
}

func getSendBackAccessTokenFromEnv() bool {
	sendAccessTokenBack := !strings.EqualFold(os.Getenv("AUTH_HEADER_RETURN"), "false")
	log.Println("AUTH_HEADER_RETURN:", sendAccessTokenBack)
	return sendAccessTokenBack
}

func getSendBackAccessTokenNameFromEnv() string {
	authHeaderName := cmp.Or(os.Getenv("SEND_ACCESS_TOKEN_HEADER_NAME"), "Authorization")
	log.Println("SEND_ACCESS_TOKEN_HEADER_NAME:", authHeaderName)
	return authHeaderName
}

func getSendAllClaimsAsJSONFromEnv() bool {
	sendAllClaimsAsJSON := !strings.EqualFold(os.Getenv("SEND_BACK_CLAIMS"), "false")
	log.Println("SEND_BACK_CLAIMS:", sendAllClaimsAsJSON)
	return sendAllClaimsAsJSON
}

func getRequiredClaimsFromEnv() []string {
	requiredClaims := os.Getenv("CLAIMS_CONTAINS")
	if requiredClaims == "" {
		log.Println("CLAIMS CONTAINS check turned off")
		return nil
	}

	requiredClaimChecks := strings.Split(requiredClaims, ",")
	log.Println("CLAIMS_CONTAINS:", requiredClaimChecks)
	for i, v := range requiredClaimChecks {
		log.Printf("\tclaim %d: value: %s\n", i, v)
	}
	return requiredClaimChecks
}

func getCacheTTLFromEnv() int {
	ttlInSeconds, err := strconv.Atoi(os.Getenv("CACHE_TTL"))
	if err != nil {
		ttlInSeconds = 300 // Default value if conversion fails
	}
	log.Println("CACHE_TTL:", ttlInSeconds)
	return ttlInSeconds
}

func getExpectedIssuerFromEnv() string {
	issuer := os.Getenv("EXPECTED_ISSUER")
	if issuer != "" {
		log.Println("EXPECTED_ISSUER:", issuer)
	} else {
		log.Println("EXPECTED_ISSUER: (not configured)")
	}
	return issuer
}

func getExpectedAudienceFromEnv() string {
	audience := os.Getenv("EXPECTED_AUDIENCE")
	if audience != "" {
		log.Println("EXPECTED_AUDIENCE:", audience)
	} else {
		log.Println("EXPECTED_AUDIENCE: (not configured)")
	}
	return audience
}
