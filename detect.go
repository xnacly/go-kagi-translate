package gokagitranslate

import "context"

type DetectParams struct {
	Text          string
	Session_token string
}

func (kt *Kagi) DetectWithParams(ctx context.Context, params DetectParams) (DetectResponse, error) {
	if err := kt.auth(ctx); err != nil {
		return DetectResponse{}, err
	}
	return DetectResponse{}, ErrNotImplemented
}

func (kt *Kagi) Detect(ctx context.Context, text string) (DetectResponse, error) {
	params := DetectParams{
		Text:          text,
		Session_token: kt.session.Token,
	}

	return kt.DetectWithParams(ctx, params)
}
