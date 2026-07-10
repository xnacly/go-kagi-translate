package gokagitranslate

import (
	"context"
	"net/http"
	"net/url"
)

var (
	kagiBase, _ = url.Parse("https://kagi.com")
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
	session   AuthResponse
}

// OneShot translates text directly using token for authentication.
func OneShot(ctx context.Context, token, from, to, text string) (string, error) {
	res, err := New(token).Translate(ctx, from, to, text)
	if err != nil {
		return "", err
	}
	return res.Translation, nil
}

// New creates a Kagi client authenticated with a Kagi session token.
//
// Using the default http.Client and DefaultUserAgent
func New(token string) *Kagi {
	return &Kagi{
		client:    &http.Client{},
		token:     token,
		userAgent: DefaultUserAgent,
	}
}

// WithClient configures the HTTP client used for requests.
func (kt *Kagi) WithClient(client *http.Client) *Kagi {
	kt.client = client
	return kt
}

// WithUserAgent configures the User-Agent header used for requests.
func (kt *Kagi) WithUserAgent(userAgent string) *Kagi {
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	kt.userAgent = userAgent
	return kt
}
