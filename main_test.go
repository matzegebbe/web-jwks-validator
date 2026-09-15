package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
