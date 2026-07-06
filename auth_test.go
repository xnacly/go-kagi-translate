package gokagitranslate

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	tok, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		t.Skip("missing KAGI_TOKEN env var")
	}

	client := New().
		WithToken(tok)

	err := client.auth(t.Context())
	if err != nil {
		t.Error("Failed to authenticate: ", err)
	}
}

func TestDecodeAuthResponseRejectsNull(t *testing.T) {
	_, err := decodeAuthResponse(strings.NewReader("null"))
	if !errors.Is(err, ErrAuthNullResponse) {
		t.Fatalf("expected ErrAuthNullResponse, got %v", err)
	}
}
