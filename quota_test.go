package gokagitranslate

import (
	"os"
	"testing"
)

func TestQuota(t *testing.T) {
	tok, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		t.Skip("missing KAGI_TOKEN env var")
	}

	client := New().WithToken(tok)
	quota, err := client.Quota(t.Context())
	if err != nil {
		t.Error("Failed to fetch quota: ", err)
	}

	if quota.Translate.Kind == "" {
		t.Error("missing translate quota kind")
	}
}
