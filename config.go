package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetPathFromEnv() string {
	path := os.Getenv("SERVER_PATH")
	if path == "" {
		fmt.Println("No SERVER_PATH environment variable detected, defaulting to /")
		path = "/"
	}
	fmt.Println("SERVER_PATH: ", path)
	return path
}
func GetPortFromEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("No PORT environment variable detected, defaulting to 8080")
		port = "8080"
	}
	fmt.Println("PORT: ", port)
	return port
}
func GetJwksURLFromEnv() string {
	jwksUrl := os.Getenv("JWKS_URL")
	if jwksUrl == "" {
		fmt.Println("No JWKS_URL environment variable detected, defaulting to " +
			"https://login.windows.net/common/discovery/keys")
		jwksUrl = "https://login.windows.net/common/discovery/keys"
	}
	fmt.Println("JWKS_URL: ", jwksUrl)
	return jwksUrl
}

func GetAuthHeaderNameFromEnv() string {
	authHeaderName := os.Getenv("AUTH_HEADER_NAME")
	if authHeaderName == "" {
		fmt.Println("No AUTH_HEADER_NAME environment variable detected, defaulting to Authorization")
		authHeaderName = "Authorization"
	}
	fmt.Println("AUTH_HEADER_NAME: ", authHeaderName)
	return authHeaderName
}

func GetAccessTokenSettingFromEnv() bool {
	sendAccessTokenBackEnv := os.Getenv("AUTH_HEADER_RETURN")
	sendAccessTokenBack := true
	if strings.ToLower(sendAccessTokenBackEnv) == "false" {
		fmt.Println("AUTH_HEADER_RETURN set to false. Will not return auth header")
		sendAccessTokenBack = false
	}
	fmt.Println("AUTH_HEADER_RETURN: ", sendAccessTokenBack)
	return sendAccessTokenBack
}

func GetSendAllClaimsAsJson() bool {
	sendAllClaimsAsJsonEnv := os.Getenv("SEND_BACK_CLAIMS")
	sendAllClaimsAsJson := true
	if strings.ToLower(sendAllClaimsAsJsonEnv) == "false" {
		fmt.Println("SEND_BACK_CLAIMS set to false. Will not return auth header")
		sendAllClaimsAsJson = false
	}
	fmt.Println("SEND_BACK_CLAIMS: ", sendAllClaimsAsJson)
	return sendAllClaimsAsJson
}

func GetTTLFromEnv() int {
	cacheTTLEnv := os.Getenv("CACHE_TTL")
	ttlInSeconds, err := strconv.Atoi(cacheTTLEnv)
	if err != nil {
		fmt.Printf("No valid CACHE_TTL environment variable detected, defaulting to " +
			"300 seconds")
		ttlInSeconds = 300 // Default value if conversion fails
	}
	fmt.Println("CACHE_TTL: ", ttlInSeconds)
	return ttlInSeconds
}
