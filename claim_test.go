package main

import (
	"net/http"
	"testing"
)

func TestClaimsContainAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		claims             map[string]any
		claimContainsCheck []string
		want               bool
	}{
		{
			name: "all required claims exist",
			claims: map[string]any{
				"name": "John",
				"tags": []any{"student", "admin"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			want:               true,
		},
		{
			name: "required slice value is missing",
			claims: map[string]any{
				"name": "John",
				"tags": []string{"student", "developer"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			want:               false,
		},
		{
			name: "required string value is missing",
			claims: map[string]any{
				"name": "Ralle",
				"tags": []string{"student", "developer"},
			},
			claimContainsCheck: []string{"name=John", "tags=admin"},
			want:               false,
		},
		{
			name: "claim value contains separator",
			claims: map[string]any{
				"scope": "resource=read",
			},
			claimContainsCheck: []string{"scope=resource=read"},
			want:               true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := claimsContainAll(tt.claims, tt.claimContainsCheck); got != tt.want {
				t.Errorf("claimsContainAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "bearer token", value: "Bearer signed-token", want: "signed-token"},
		{name: "raw token", value: "signed-token", want: "signed-token"},
		{name: "repeated whitespace", value: "Bearer   signed-token", want: "signed-token"},
		{name: "empty header", want: ""},
		{name: "too many fields", value: "Bearer signed-token unexpected", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatalf("http.NewRequest() error = %v", err)
			}
			req.Header.Set("Authorization", tt.value)

			if got := extractToken(req, "Authorization"); got != tt.want {
				t.Errorf("extractToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
