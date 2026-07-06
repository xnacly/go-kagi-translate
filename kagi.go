package gokagitranslate

import (
	"context"
	"net/http"
)

const (
	translate = "https://translate.kagi.com/api/translate"
	auth      = "https://translate.kagi.com/api/auth"
	quota     = "https://translate.kagi.com/api/quota"

	// Most common user agent string
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.10 Safari/605.1.1"
)

type Kagi struct {
	ctx     context.Context
	client  *http.Client
	token   string
	session string
}

func New() *Kagi {
	return &Kagi{
		ctx:    context.Background(),
		client: http.DefaultClient,
	}
}

func (kt *Kagi) WithCtx(ctx context.Context) *Kagi {
	kt.ctx = ctx
	return kt
}

func (kt *Kagi) WithClient(client *http.Client) *Kagi {
	kt.client = client
	return kt
}

func (kt *Kagi) WithToken(token string) *Kagi {
	kt.token = token
	return kt
}
