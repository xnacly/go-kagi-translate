package gokagitranslate

import (
	"net/http"
	"net/http/cookiejar"
)

func (kt *Kagi) prepReq(req *http.Request) error {
	if kt.client.Jar == nil {
		kt.client.Jar, _ = cookiejar.New(nil)
	}

	req.Header.Set("User-Agent", kt.userAgent)

	// TODO: support v1 request signing once the API contract is stable.

	kt.client.Jar.SetCookies(kagiBase, []*http.Cookie{
		{Name: "kagi_session", Value: kt.token, Domain: ".kagi.com"},
	})

	return nil
}
