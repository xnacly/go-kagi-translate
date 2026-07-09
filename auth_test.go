package gokagitranslate

import (
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	tok, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		t.Skip("missing KAGI_TOKEN env var")
	}

	client := New(tok)

	err := client.auth(t.Context())
	if err != nil {
		t.Error("Failed to authenticate: ", err)
	}
}
