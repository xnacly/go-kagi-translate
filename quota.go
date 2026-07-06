package gokagitranslate

import (
	"context"
	"fmt"
	"net/http"
)

func (kt *Kagi) Quota(ctx context.Context) (QuotaResponse, error) {
	if err := kt.auth(ctx); err != nil {
		return QuotaResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", quota, nil)
	if err != nil {
		return QuotaResponse{}, err
	}

	kt.prepReq(req)

	res, err := kt.client.Do(req)
	if err != nil {
		return QuotaResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return QuotaResponse{}, fmt.Errorf("quota failed: %s", res.Status)
	}

	return decodeResponse[QuotaResponse](res.Body)
}
