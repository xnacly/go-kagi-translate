package gokagitranslate

import (
	"net/http"
)

func (kt *Kagi) prepReq(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Referer", "https://translate.kagi.com/")
	req.Header.Set("User-Agent", kt.userAgent)

	req.AddCookie(&http.Cookie{Name: "kagi_session", Value: kt.token})
}
