package gokagitranslate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrAuthNullResponse = errors.New("auth failed: empty session response")

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

	session, err := decodeAuthResponse(res.Body)
	if err != nil {
		return err
	}

	kt.session = session
	return nil
}

func decodeAuthResponse(body io.Reader) (AuthResponse, error) {
	d := json.NewDecoder(body)
	var session *AuthResponse
	if err := d.Decode(&session); err != nil {
		return AuthResponse{}, err
	}
	if session == nil {
		return AuthResponse{}, ErrAuthNullResponse
	}

	return *session, nil
}
