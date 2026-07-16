package gokagitranslate

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const authExpirySkew = 30 * time.Second

func (kt *Kagi) auth(ctx context.Context) (AuthResponse, error) {
	if session, ok := kt.cachedAuth(time.Now()); ok {
		return session, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", auth, nil)
	if err != nil {
		return AuthResponse{}, err
	}

	kt.prepReq(req)

	res, err := kt.client.Do(req)
	if err != nil {
		return AuthResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return AuthResponse{}, fmt.Errorf("auth failed: %s", res.Status)
	}

	session, err := decodeResponse[AuthResponse](res.Body)
	if err != nil {
		return AuthResponse{}, err
	}
	if session.Token == "" {
		return AuthResponse{}, fmt.Errorf("auth failed: empty session token")
	}
	if !session.ExpiresAt.IsZero() && time.Now().After(session.ExpiresAt) {
		return AuthResponse{}, fmt.Errorf("auth failed: expired session token")
	}
	kt.storeAuth(session, time.Now())

	return session, nil
}

func (kt *Kagi) cachedAuth(now time.Time) (AuthResponse, bool) {
	session := kt.authCache.Load()
	if session == nil || !usableAuth(*session, now) {
		return AuthResponse{}, false
	}
	return *session, true
}

func (kt *Kagi) storeAuth(session AuthResponse, now time.Time) {
	if !cacheableAuth(session, now) {
		return
	}
	candidate := session
	for {
		current := kt.authCache.Load()
		if current != nil && usableAuth(*current, now) && current.ExpiresAt.After(candidate.ExpiresAt) {
			return
		}
		if kt.authCache.CompareAndSwap(current, &candidate) {
			return
		}
	}
}

func usableAuth(session AuthResponse, now time.Time) bool {
	return session.Token != "" && cacheableAuth(session, now)
}

func cacheableAuth(session AuthResponse, now time.Time) bool {
	return !session.ExpiresAt.IsZero() && now.Add(authExpirySkew).Before(session.ExpiresAt)
}
