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
	auth      = "https://translate.kagi.com/api/auth"
	quota     = "https://translate.kagi.com/api/quota"

	// We are being nice and adhere to the netiquette, this way kagi can send me an email if this is being abused, i guess?
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.10 Safari/605.1.1 go-kagi-translate contact@xnacly.me"
)

type Kagi struct {
	client    *http.Client
	token     string
	userAgent string
	session   AuthResponse
}

// OneShot skips the builder style api and translates "text" directly "from" to
// "to" using "token" to authenticate
func OneShot(ctx context.Context, token, from, to, text string) (string, error) {
	// TODO: translate here
	return "", nil
}

func New(token string) *Kagi {
	return &Kagi{
		client:    &http.Client{},
		token:     token,
		userAgent: DefaultUserAgent,
	}
}

func (kt *Kagi) WithClient(client *http.Client) *Kagi {
	kt.client = client
	return kt
}

func (kt *Kagi) WithToken(token string) *Kagi {
	kt.token = token
	return kt
}

func (kt *Kagi) WithUserAgent(userAgent string) *Kagi {
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	kt.userAgent = userAgent
	return kt
}
