package gokagitranslate

import (
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	_ = New().
		WithClient(&http.Client{}).
		WithToken("12345")
}
