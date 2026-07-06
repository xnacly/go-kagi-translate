package gokagitranslate

import (
	"errors"
	"strings"
	"testing"
)

func TestDecodeResponseRejectsNull(t *testing.T) {
	_, err := decodeResponse[AuthResponse](strings.NewReader("null"))
	if !errors.Is(err, ErrNullResponse) {
		t.Fatalf("expected ErrAuthNullResponse, got %v", err)
	}
}
