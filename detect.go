package gokagitranslate

import (
	"context"
	"encoding/json"
	"io"
)

// DetectParams contains the request body sent to Kagi's detect endpoint.
type DetectParams struct {
	Text                string   `json:"text"`
	IncludeAlternatives bool     `json:"include_alternatives"`
	RecentLanguages     []string `json:"recent_languages,omitempty"`
	SessionToken        string   `json:"session_token,omitempty"`
}

// DetectWithParams detects the language of text using the provided options.
func (kt *Kagi) DetectWithParams(ctx context.Context, params DetectParams) (DetectResponse, error) {
	session, err := kt.auth(ctx)
	if err != nil {
		return DetectResponse{}, err
	}
	if params.SessionToken == "" {
		params.SessionToken = session.Token
	}

	res, err := kt.doJSON(ctx, "POST", detect, "detect", params)
	if err != nil {
		return DetectResponse{}, err
	}
	defer res.Body.Close()

	return decodeDetectResponse(res.Body)
}

// Detect detects the language of text using standard options.
func (kt *Kagi) Detect(ctx context.Context, text string) (DetectResponse, error) {
	params := DetectParams{
		Text:                text,
		IncludeAlternatives: true,
	}

	return kt.DetectWithParams(ctx, params)
}

func decodeDetectResponse(body io.Reader) (DetectResponse, error) {
	var raw struct {
		DetectedLanguage Language   `json:"detected_language"`
		Language         Language   `json:"language"`
		Iso              string     `json:"iso"`
		Label            string     `json:"label"`
		Alternatives     []Language `json:"alternatives"`
	}
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return DetectResponse{}, err
	}

	out := DetectResponse{
		DetectedLanguage: raw.DetectedLanguage,
		Alternatives:     raw.Alternatives,
	}
	if out.DetectedLanguage.Iso == "" && out.DetectedLanguage.Label == "" {
		out.DetectedLanguage = raw.Language
	}
	if out.DetectedLanguage.Iso == "" && out.DetectedLanguage.Label == "" {
		out.DetectedLanguage = Language{Iso: raw.Iso, Label: raw.Label}
	}
	if out.DetectedLanguage.Iso == "" && out.DetectedLanguage.Label == "" {
		return DetectResponse{}, ErrEmptyResponse
	}

	return out, nil
}
