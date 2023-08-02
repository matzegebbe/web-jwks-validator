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
		path = "/"
	}
	fmt.Println("SERVER_PATH: ", path)
	return path
}
func GetPortFromEnv() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("PORT: ", port)
	return port
}
func GetJwksURLFromEnv() string {
	jwksUrl := os.Getenv("JWKS_URL")
	if jwksUrl == "" {
		jwksUrl = "https://login.windows.net/common/discovery/keys"
	}
	fmt.Println("JWKS_URL: ", jwksUrl)
	return jwksUrl
}

func GetAuthHeaderNameFromEnv() string {
	authHeaderName := os.Getenv("AUTH_HEADER_NAME")
	if authHeaderName == "" {
		authHeaderName = "Authorization"
	}
	fmt.Println("AUTH_HEADER_NAME: ", authHeaderName)
	return authHeaderName
}

func GetAccessTokenSettingFromEnv() bool {
	sendAccessTokenBackEnv := os.Getenv("AUTH_HEADER_RETURN")
	sendAccessTokenBack := true
	if strings.ToLower(sendAccessTokenBackEnv) == "false" {
		sendAccessTokenBack = false
	}
	fmt.Println("AUTH_HEADER_RETURN: ", sendAccessTokenBack)
	return sendAccessTokenBack
}

func GetSendAllClaimsAsJson() bool {
	sendAllClaimsAsJsonEnv := os.Getenv("SEND_BACK_CLAIMS")
	sendAllClaimsAsJson := true
	if strings.ToLower(sendAllClaimsAsJsonEnv) == "false" {
		sendAllClaimsAsJson = false
	}
	fmt.Println("SEND_BACK_CLAIMS: ", sendAllClaimsAsJson)
	return sendAllClaimsAsJson
}

func GetTTLFromEnv() int {
	cacheTTLEnv := os.Getenv("CACHE_TTL")
	ttlInSeconds, err := strconv.Atoi(cacheTTLEnv)
	if err != nil {
		ttlInSeconds = 300 // Default value if conversion fails
	}
	fmt.Println("CACHE_TTL: ", ttlInSeconds)
	return ttlInSeconds
}
