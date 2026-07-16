package gokagitranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestTranslatePostsReferencePayload(t *testing.T) {
	var requests []*http.Request
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests = append(requests, req)
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"translate-session-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}

			if req.URL.String() != translate {
				t.Fatalf("got url %q, want %q", req.URL.String(), translate)
			}
			if req.Method != http.MethodPost {
				t.Fatalf("got method %q, want %q", req.Method, http.MethodPost)
			}
			if got := req.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("got content type %q, want application/json", got)
			}
			if got := req.Header.Get("X-Signal"); got != "abortable" {
				t.Fatalf("got X-Signal %q, want abortable", got)
			}

			assertTranslatePayload(t, req, "me llamo matteo", "es", "en", "translate-session-token")
			return jsonResponse(http.StatusOK, `{"translation":"My name is Matteo","detected_language":{"iso":"es","label":"Spanish"}}`), nil
		}),
	}

	res, err := New("kagi-session-token", WithClient(client)).Translate(t.Context(), "es", "en", "me llamo matteo")
	if err != nil {
		t.Fatal(err)
	}
	if res.Translation != "My name is Matteo" {
		t.Fatalf("got translation %q, want %q", res.Translation, "My name is Matteo")
	}
	if len(requests) != 2 {
		t.Fatalf("got %d requests, want auth and translate", len(requests))
	}
}

func TestTranslateWithParamsUsesProvidedSessionToken(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"auth-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			var body struct {
				SessionToken string `json:"session_token"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.SessionToken != "provided-token" {
				t.Fatalf("got session token %q, want provided-token", body.SessionToken)
			}
			return jsonResponse(http.StatusOK, `{"translation":"hola"}`), nil
		}),
	}

	_, err := New("kagi-session-token", WithClient(client)).TranslateWithParams(t.Context(), TranslateParams{
		Text:         "hello",
		From:         "en",
		To:           "es",
		SessionToken: "provided-token",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTranslateFailsOnNonSuccessStatus(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"auth-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			return jsonResponse(http.StatusForbidden, `{"error":"forbidden"}`), nil
		}),
	}

	_, err := New("kagi-session-token", WithClient(client)).Translate(t.Context(), "en", "es", "hello")
	if err == nil || !strings.Contains(err.Error(), "translate failed: 403 Forbidden") {
		t.Fatalf("got error %v, want translate failed status", err)
	}
}

func TestTranslateReusesCachedAuth(t *testing.T) {
	var authRequests atomic.Int64
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				authRequests.Add(1)
				return jsonResponse(http.StatusOK, `{"token":"fresh-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			return jsonResponse(http.StatusOK, `{"translation":"hola"}`), nil
		}),
	}
	kt := New("kagi-session-token", WithClient(client))

	for range 2 {
		res, err := kt.Translate(t.Context(), "en", "es", "hello")
		if err != nil {
			t.Fatal(err)
		}
		if res.Translation != "hola" {
			t.Fatalf("got translation %q, want hola", res.Translation)
		}
	}
	if got := authRequests.Load(); got != 1 {
		t.Fatalf("got %d auth requests, want 1", got)
	}
}

func TestTranslateRefreshesExpiredCachedAuth(t *testing.T) {
	var authRequests atomic.Int64
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				authRequests.Add(1)
				return jsonResponse(http.StatusOK, `{"token":"fresh-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			return jsonResponse(http.StatusOK, `{"translation":"hola"}`), nil
		}),
	}
	kt := New("kagi-session-token", WithClient(client))
	kt.authCache.Store(&AuthResponse{Token: "expired-token", ExpiresAt: time.Now().Add(-time.Minute)})

	res, err := kt.Translate(t.Context(), "en", "es", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if res.Translation != "hola" {
		t.Fatalf("got translation %q, want hola", res.Translation)
	}
	if got := authRequests.Load(); got != 1 {
		t.Fatalf("got %d auth requests, want 1", got)
	}
}

func TestTranslateConcurrentCallsDoNotShareSessionState(t *testing.T) {
	var authRequests atomic.Int64
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				n := authRequests.Add(1)
				return jsonResponse(http.StatusOK, fmt.Sprintf(`{"token":"session-%d","expiresAt":"2030-01-01T00:00:00Z"}`, n)), nil
			}

			var body struct {
				SessionToken string `json:"session_token"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, err
			}
			if body.SessionToken == "" {
				return nil, errors.New("missing session token")
			}
			return jsonResponse(http.StatusOK, `{"translation":"hola"}`), nil
		}),
	}
	kt := New("kagi-session-token", WithClient(client))

	const calls = 16
	var wg sync.WaitGroup
	errs := make(chan error, calls)
	for range calls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := kt.Translate(t.Context(), "en", "es", "hello")
			if err != nil {
				errs <- err
				return
			}
			if res.Translation != "hola" {
				errs <- fmt.Errorf("got translation %q, want hola", res.Translation)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := authRequests.Load(); got == 0 || got > calls {
		t.Fatalf("got %d auth requests, want between 1 and %d", got, calls)
	}
}

func TestDoJSONSetsJSONHeadersForPayload(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got := req.Header.Get("Content-Type"); got != "application/json" {
				t.Fatalf("got content type %q, want application/json", got)
			}
			if got := req.Header.Get("Origin"); got != "https://translate.kagi.com" {
				t.Fatalf("got origin %q, want https://translate.kagi.com", got)
			}
			if got := req.Header.Get("X-Signal"); got != "abortable" {
				t.Fatalf("got X-Signal %q, want abortable", got)
			}
			return jsonResponse(http.StatusOK, `{"translation":"hola"}`), nil
		}),
	}

	res, err := New("kagi-session-token", WithClient(client)).doJSON(t.Context(), http.MethodPost, translate, "translate", map[string]string{"text": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
}

func assertTranslatePayload(t *testing.T, req *http.Request, text, from, to, sessionToken string) {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"text":                     text,
		"from":                     from,
		"to":                       to,
		"stream":                   false,
		"formality":                "default",
		"speaker_gender":           "unknown",
		"addressee_gender":         "unknown",
		"language_complexity":      "standard",
		"translation_style":        "natural",
		"context":                  "",
		"model":                    "standard",
		"session_token":            sessionToken,
		"dictionary_language":      to,
		"use_definition_context":   false,
		"enable_language_features": false,
	}
	for key, wantValue := range want {
		if body[key] != wantValue {
			t.Fatalf("payload[%s] = %#v, want %#v", key, body[key], wantValue)
		}
	}
	if _, ok := body["predicted_language"]; ok {
		t.Fatal("payload unexpectedly included empty predicted_language")
	}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
