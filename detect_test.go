package gokagitranslate

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

func TestDetectPostsReferencePayload(t *testing.T) {
	var requests []*http.Request
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests = append(requests, req)
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"detect-session-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}

			if req.URL.String() != detect {
				t.Fatalf("got url %q, want %q", req.URL.String(), detect)
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

			assertDetectPayload(t, req, map[string]any{
				"text":                 "skibid",
				"include_alternatives": true,
				"session_token":        "detect-session-token",
			}, "recent_languages")
			return jsonResponse(http.StatusOK, `{"detected_language":{"iso":"en","label":"English"}}`), nil
		}),
	}

	res, err := New("kagi-session-token", WithClient(client)).Detect(t.Context(), "skibid")
	if err != nil {
		t.Fatal(err)
	}
	if res.DetectedLanguage.Iso != "en" {
		t.Fatalf("got language %q, want en", res.DetectedLanguage.Iso)
	}
	if len(requests) != 2 {
		t.Fatalf("got %d requests, want auth and detect", len(requests))
	}
}

func TestDetectWithParamsUsesProvidedSessionToken(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"auth-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			assertDetectPayload(t, req, map[string]any{
				"text":                 "hola",
				"include_alternatives": true,
				"recent_languages":     []any{"es", "en"},
				"session_token":        "provided-token",
			})
			return jsonResponse(http.StatusOK, `{"detected_language":{"iso":"es","label":"Spanish"}}`), nil
		}),
	}

	_, err := New("kagi-session-token", WithClient(client)).DetectWithParams(t.Context(), DetectParams{
		Text:                "hola",
		IncludeAlternatives: true,
		RecentLanguages:     []string{"es", "en"},
		SessionToken:        "provided-token",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDetectFailsOnNonSuccessStatus(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				return jsonResponse(http.StatusOK, `{"token":"auth-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			return jsonResponse(http.StatusForbidden, `{"error":"forbidden"}`), nil
		}),
	}

	_, err := New("kagi-session-token", WithClient(client)).Detect(t.Context(), "hello")
	if err == nil || !strings.Contains(err.Error(), "detect failed: 403 Forbidden") {
		t.Fatalf("got error %v, want detect failed status", err)
	}
}

func TestDetectDecodesAlternateResponseShapes(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "detected_language", body: `{"detected_language":{"iso":"en","label":"English"}}`},
		{name: "language", body: `{"language":{"iso":"en","label":"English"}}`},
		{name: "top_level", body: `{"iso":"en","label":"English"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := decodeDetectResponse(strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			if res.DetectedLanguage.Iso != "en" {
				t.Fatalf("got language %q, want en", res.DetectedLanguage.Iso)
			}
		})
	}
}

func TestDetectReusesCachedAuth(t *testing.T) {
	var authRequests atomic.Int64
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() == auth {
				authRequests.Add(1)
				return jsonResponse(http.StatusOK, `{"token":"fresh-token","expiresAt":"2030-01-01T00:00:00Z"}`), nil
			}
			return jsonResponse(http.StatusOK, `{"detected_language":{"iso":"en","label":"English"}}`), nil
		}),
	}
	kt := New("kagi-session-token", WithClient(client))

	for range 2 {
		res, err := kt.Detect(t.Context(), "hello")
		if err != nil {
			t.Fatal(err)
		}
		if res.DetectedLanguage.Iso != "en" {
			t.Fatalf("got language %q, want en", res.DetectedLanguage.Iso)
		}
	}
	if got := authRequests.Load(); got != 1 {
		t.Fatalf("got %d auth requests, want 1", got)
	}
}

func assertDetectPayload(t *testing.T, req *http.Request, want map[string]any, omitted ...string) {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	for key, wantValue := range want {
		if !equalJSONValue(body[key], wantValue) {
			t.Fatalf("payload[%s] = %#v, want %#v", key, body[key], wantValue)
		}
	}
	for _, key := range omitted {
		if _, ok := body[key]; ok {
			t.Fatalf("payload unexpectedly included %q", key)
		}
	}
}

func equalJSONValue(got, want any) bool {
	gotBytes, err := json.Marshal(got)
	if err != nil {
		return false
	}
	wantBytes, err := json.Marshal(want)
	if err != nil {
		return false
	}
	return string(gotBytes) == string(wantBytes)
}
