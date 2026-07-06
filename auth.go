package gokagitranslate

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (kt *Kagi) auth(ctx context.Context) error {
	if time.Now().Before(kt.session.ExpiresAt) {
		// we skip reauth if session is still valid
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", auth, nil)
	if err != nil {
		return err
	}

	kt.prepReq(req)

	res, err := kt.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("auth failed: %s", res.Status)
	}

	session, err := decodeResponse[AuthResponse](res.Body)
	if err != nil {
		return err
	}

	kt.session = session
	return nil
}
