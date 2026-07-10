package gokagitranslate

import (
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	_ = New("12345").
		WithClient(&http.Client{}).
		WithUserAgent("custom-user-agent")
}

func TestNewSetsToken(t *testing.T) {
	client := New("12345")
	if client.token != "12345" {
		t.Fatalf("got token %q, want %q", client.token, "12345")
	}
}

func TestPrepReqUsesDefaultUserAgent(t *testing.T) {
	req, err := http.NewRequest("GET", auth, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := New("12345").prepReq(req); err != nil {
		t.Fatal(err)
	}

	if got := req.Header.Get("User-Agent"); got != DefaultUserAgent {
		t.Fatalf("got user agent %q, want %q", got, DefaultUserAgent)
	}
}

func TestPrepReqUsesConfiguredUserAgent(t *testing.T) {
	const userAgent = "custom-user-agent"
	req, err := http.NewRequest("GET", auth, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := New("12345").WithUserAgent(userAgent).prepReq(req); err != nil {
		t.Fatal(err)
	}

	if got := req.Header.Get("User-Agent"); got != userAgent {
		t.Fatalf("got user agent %q, want %q", got, userAgent)
	}
}

func TestPrepReqDoesNotSetJSONHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", auth, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := New("12345").prepReq(req); err != nil {
		t.Fatal(err)
	}

	if got := req.Header.Get("Content-Type"); got != "" {
		t.Fatalf("got content type %q, want empty", got)
	}
	if got := req.Header.Get("X-Signal"); got != "" {
		t.Fatalf("got X-Signal %q, want empty", got)
	}
}
