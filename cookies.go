package gokagitranslate

import (
	"net/http"
	"net/http/cookiejar"
)

func (kt *Kagi) prepReq(req *http.Request) error {
	if kt.client.Jar == nil {
		kt.client.Jar, _ = cookiejar.New(nil)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Referer", "https://translate.kagi.com/")
	req.Header.Set("User-Agent", kt.userAgent)

	kt.client.Jar.SetCookies(kagiBase, []*http.Cookie{
		{Name: "kagi_session", Value: kt.token, Domain: ".kagi.com"},
	})

	return nil
}
