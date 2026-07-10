package gokagitranslate

import (
	"errors"
	"strings"
	"testing"
)

func TestDecodeResponseRejectsNull(t *testing.T) {
	_, err := decodeResponse[AuthResponse](strings.NewReader("null"))
	if !errors.Is(err, ErrEmptyResponse) {
		t.Fatalf("expected ErrEmptyResponse, got %v", err)
	}
}
