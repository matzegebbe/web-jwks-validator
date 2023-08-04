package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetPathFromEnv() string {
	path := os.Getenv("SERVER_PATH")
	if path == "" {
		path = "/"
	}
	log.Println("SERVER_PATH:", path)
	return path
}
func GetPortFromEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("PORT:", port)
	return port
}
func GetJwksURLFromEnv() string {
	jwksUrl := os.Getenv("JWKS_URL")
	if jwksUrl == "" {
		jwksUrl = "https://login.windows.net/common/discovery/keys"
	}
	log.Println("JWKS_URL:", jwksUrl)
	return jwksUrl
}

func GetAuthHeaderNameFromEnv() string {
	authHeaderName := os.Getenv("AUTH_HEADER_NAME")
	if authHeaderName == "" {
		authHeaderName = "Authorization"
	}
	log.Println("AUTH_HEADER_NAME:", authHeaderName)
	return authHeaderName
}

func GetAccessTokenSettingFromEnv() bool {
	sendAccessTokenBackEnv := os.Getenv("AUTH_HEADER_RETURN")
	sendAccessTokenBack := true
	if strings.ToLower(sendAccessTokenBackEnv) == "false" {
		sendAccessTokenBack = false
	}
	log.Println("AUTH_HEADER_RETURN:", sendAccessTokenBack)
	return sendAccessTokenBack
}

func GetSendAllClaimsAsJson() bool {
	sendAllClaimsAsJsonEnv := os.Getenv("SEND_BACK_CLAIMS")
	sendAllClaimsAsJson := true
	if strings.ToLower(sendAllClaimsAsJsonEnv) == "false" {
		sendAllClaimsAsJson = false
	}
	log.Println("SEND_BACK_CLAIMS:", sendAllClaimsAsJson)
	return sendAllClaimsAsJson
}

func GetClaimContains() []string {
	claimContainsArr := strings.Split(os.Getenv("CLAIMS_CONTAINS"), ",")
	if len(claimContainsArr) != 0 {
		log.Println("CLAIMS_CONTAINS:", claimContainsArr)
		for i, v := range claimContainsArr {
			log.Printf("\tclaim %d: value: %s\n", i, v)
		}
	} else {
		log.Println("CLAIMS_CONTAINS NOT CHECKED")
	}
	return claimContainsArr
}

func GetTTLFromEnv() int {
	cacheTTLEnv := os.Getenv("CACHE_TTL")
	ttlInSeconds, err := strconv.Atoi(cacheTTLEnv)
	if err != nil {
		ttlInSeconds = 300 // Default value if conversion fails
	}
	log.Println("CACHE_TTL:", ttlInSeconds)
	return ttlInSeconds
}
