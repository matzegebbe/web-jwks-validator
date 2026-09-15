package main

import (
	"slices"
	"testing"
)

func TestStringConfigDefaults(t *testing.T) {
	tests := []struct {
		name string
		env  string
		get  func() string
		want string
	}{
		{name: "server path", env: "SERVER_PATH", get: getPathFromEnv, want: "/"},
		{name: "port", env: "PORT", get: getPortFromEnv, want: "8080"},
		{name: "JWKS URL", env: "JWKS_URL", get: getJWKSURLFromEnv, want: "https://login.windows.net/common/discovery/keys"},
		{name: "auth header", env: "AUTH_HEADER_NAME", get: getAuthHeaderNameFromEnv, want: "Authorization"},
		{name: "returned token header", env: "SEND_ACCESS_TOKEN_HEADER_NAME", get: getSendBackAccessTokenNameFromEnv, want: "Authorization"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.env, "")

			if got := tt.get(); got != tt.want {
				t.Errorf("config value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBooleanConfig(t *testing.T) {
	tests := []struct {
		name string
		env  string
		get  func() bool
	}{
		{name: "return access token", env: "AUTH_HEADER_RETURN", get: getSendBackAccessTokenFromEnv},
		{name: "return claims", env: "SEND_BACK_CLAIMS", get: getSendAllClaimsAsJSONFromEnv},
	}

	for _, tt := range tests {
		t.Run(tt.name+" defaults to true", func(t *testing.T) {
			t.Setenv(tt.env, "")
			if !tt.get() {
				t.Error("config value = false, want true")
			}
		})
		t.Run(tt.name+" accepts case-insensitive false", func(t *testing.T) {
			t.Setenv(tt.env, "FaLsE")
			if tt.get() {
				t.Error("config value = true, want false")
			}
		})
	}
}

func TestGetRequiredClaimsFromEnv(t *testing.T) {
	t.Setenv("CLAIMS_CONTAINS", "roles=writer,scope=documents=read")

	want := []string{"roles=writer", "scope=documents=read"}
	if got := getRequiredClaimsFromEnv(); !slices.Equal(got, want) {
		t.Errorf("getRequiredClaimsFromEnv() = %v, want %v", got, want)
	}
}

func TestGetCacheTTLFromEnv(t *testing.T) {
	t.Run("configured", func(t *testing.T) {
		t.Setenv("CACHE_TTL", "60")
		if got := getCacheTTLFromEnv(); got != 60 {
			t.Errorf("getCacheTTLFromEnv() = %d, want 60", got)
		}
	})

	t.Run("invalid uses default", func(t *testing.T) {
		t.Setenv("CACHE_TTL", "invalid")
		if got := getCacheTTLFromEnv(); got != 300 {
			t.Errorf("getCacheTTLFromEnv() = %d, want 300", got)
		}
	})
}
