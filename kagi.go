package gokagitranslate

import (
	"context"
	"net/http"
)

const (
	translate = "https://translate.kagi.com/api/translate"
	auth      = "https://translate.kagi.com/api/auth"
	quota     = "https://translate.kagi.com/api/quota"
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:151.0) Gecko/20100101 Firefox/151.0 go-kagi-translate"
)

type Kagi struct {
	ctx     context.Context
	client  http.Client
	token   string
	session string
}

func New() Kagi { return Kagi{} }
func (kt Kagi) WithCtx(ctx context.Context) Kagi {
	kt.ctx = ctx
	return kt
}
func (kt Kagi) WithClient(client http.Client) Kagi {
	kt.client = client
	return kt
}
func (kt Kagi) WithToken(token string) Kagi {
	kt.token = token
	return kt
}
