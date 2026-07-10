package gokagitranslate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxErrorBody = 1024

func (kt *Kagi) doJSON(ctx context.Context, method, url, failure string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if err := kt.prepReq(req); err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://translate.kagi.com")
		req.Header.Set("X-Signal", "abortable")
	}

	res, err := kt.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		defer res.Body.Close()
		detail, _ := io.ReadAll(io.LimitReader(res.Body, maxErrorBody))
		if msg := strings.TrimSpace(string(detail)); msg != "" {
			return nil, fmt.Errorf("%s failed: %s: %s", failure, res.Status, msg)
		}
		return nil, fmt.Errorf("%s failed: %s", failure, res.Status)
	}

	return res, nil
}

func decodeJSON[T any](ctx context.Context, kt *Kagi, method, url, failure string, payload any) (T, error) {
	res, err := kt.doJSON(ctx, method, url, failure, payload)
	var zero T
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	return decodeResponse[T](res.Body)
}
