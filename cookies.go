package gokagitranslate

import (
	"net/http"
	"net/http/cookiejar"
)

func (kt *Kagi) prepReq(req *http.Request) {
	if kt.client.Jar == nil {
		kt.client.Jar, _ = cookiejar.New(nil)
	}

	req.Header.Add("User-Agent", UserAgent)

	kt.client.Jar.SetCookies(kagiBase, []*http.Cookie{
		{Name: "kagi_session", Value: kt.token, Domain: ".kagi.com"},
	})
}
