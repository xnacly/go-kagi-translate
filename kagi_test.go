package gokagitranslate

import (
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	_ = New().
		WithClient(&http.Client{}).
		WithCtx(t.Context()).
		WithToken("12345")
}
