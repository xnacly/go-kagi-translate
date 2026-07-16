package gokagitranslate

import (
	"context"
	"net/http"
	"sync/atomic"
)

const (
	base      = "https://translate.kagi.com/"
	translate = "https://translate.kagi.com/api/translate"
	detect    = "https://translate.kagi.com/api/detect"
	auth      = "https://translate.kagi.com/api/auth"
	quota     = "https://translate.kagi.com/api/quota"

	// DefaultUserAgent identifies this unofficial client to Kagi.
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.10 Safari/605.1.1 go-kagi-translate contact@xnacly.me"
)

// Kagi is a client for Kagi Translate's private web API.
type Kagi struct {
	client    *http.Client
	token     string
	userAgent string
	authCache atomic.Pointer[AuthResponse]
}

// Option configures a Kagi client at construction time.
type Option func(*Kagi)

// OneShot translates text directly using token for authentication.
func OneShot(ctx context.Context, token, from, to, text string) (string, error) {
	res, err := New(token).Translate(ctx, from, to, text)
	if err != nil {
		return "", err
	}
	return res.Translation, nil
}

// New creates a Kagi client authenticated with a Kagi session token.
func New(token string, options ...Option) *Kagi {
	kt := &Kagi{
		client:    &http.Client{},
		token:     token,
		userAgent: DefaultUserAgent,
	}
	for _, option := range options {
		option(kt)
	}
	return kt
}

// WithClient configures the HTTP client used for requests.
func WithClient(client *http.Client) Option {
	return func(kt *Kagi) {
		if client != nil {
			kt.client = client
		}
	}
}

// WithUserAgent configures the User-Agent header used for requests.
func WithUserAgent(userAgent string) Option {
	return func(kt *Kagi) {
		if userAgent == "" {
			userAgent = DefaultUserAgent
		}
		kt.userAgent = userAgent
	}
}
