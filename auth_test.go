package gokagitranslate

import (
	"os"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	tok, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		t.Skip("missing KAGI_TOKEN env var")
	}

	client := New(tok)

	_, err := client.auth(t.Context())
	if err != nil {
		t.Error("Failed to authenticate: ", err)
	}
}

func TestAuthCacheRejectsNearExpirySession(t *testing.T) {
	client := New("12345")
	now := time.Now()
	client.authCache.Store(&AuthResponse{
		Token:     "near-expiry-token",
		ExpiresAt: now.Add(authExpirySkew / 2),
	})

	if session, ok := client.cachedAuth(now); ok {
		t.Fatalf("got cached session %q, want cache miss", session.Token)
	}
}
