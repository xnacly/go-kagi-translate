package gokagitranslate

import "context"

// Quota returns the authenticated account's Kagi Translate quota usage.
func (kt *Kagi) Quota(ctx context.Context) (QuotaResponse, error) {
	if err := kt.auth(ctx); err != nil {
		return QuotaResponse{}, err
	}

	return decodeJSON[QuotaResponse](ctx, kt, "GET", quota, "quota", nil)
}
