package gokagitranslate

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type TranslateResponse struct {
	Translation      string `json:"translation"`
	DetectedLanguage struct {
		Iso   string `json:"iso"`
		Label string `json:"label"`
	} `json:"detected_language"`
	Definition struct {
		Word           string `json:"word"`
		PrimaryMeaning struct {
			Definition            string   `json:"definition"`
			PartOfSpeech          []string `json:"part_of_speech"`
			PartOfSpeechCanonical []string `json:"part_of_speech_canonical"`
			UsageLevel            []string `json:"usage_level"`
			Synonyms              []string `json:"synonyms"`
			SynonymComparisons    []struct {
				Synonym    string `json:"synonym"`
				Difference string `json:"difference"`
			} `json:"synonym_comparisons,omitempty"`
		} `json:"primary_meaning"`
		SecondaryMeanings []struct {
			Definition            string   `json:"definition"`
			PartOfSpeech          []string `json:"part_of_speech"`
			PartOfSpeechCanonical []string `json:"part_of_speech_canonical"`
			UsageLevel            []string `json:"usage_level"`
			Synonyms              []string `json:"synonyms"`
			SynonymComparisons    []struct {
				Synonym    string `json:"synonym"`
				Difference string `json:"difference"`
			} `json:"synonym_comparisons,omitempty"`
		} `json:"secondary_meanings"`
		Examples      []string `json:"examples"`
		Pronunciation string   `json:"pronunciation"`
		Etymology     string   `json:"etymology"`
		Notes         string   `json:"notes"`
		TemporalTrend string   `json:"temporal_trend"`
		RelatedWords  []struct {
			Word         string `json:"word"`
			Relationship string `json:"relationship"`
		} `json:"related_words"`
	} `json:"definition"`
}

type DetectResponse struct {
	DetectedLanguage struct {
		Iso   string `json:"iso"`
		Label string `json:"label"`
	} `json:"detected_language"`
}

type AuthResponse struct {
	Token              string    `json:"token"`
	ID                 string    `json:"id"`
	LoggedIn           bool      `json:"loggedIn"`
	Subscription       bool      `json:"subscription"`
	ExpiresAt          time.Time `json:"expiresAt"`
	Theme              string    `json:"theme"`
	MobileTheme        string    `json:"mobileTheme"`
	CustomCSSEnabled   bool      `json:"customCssEnabled"`
	Language           string    `json:"language"`
	CustomCSSAvailable bool      `json:"customCssAvailable"`
	AccountType        string    `json:"accountType"`
	Platform           string    `json:"platform"`
}

type QuotaResponse struct {
	Translate Quota `json:"translate"`
	Proofread Quota `json:"proofread"`
	Document  Quota `json:"document"`
}

type Quota struct {
	Kind        string    `json:"kind"`
	Used        int       `json:"used"`
	Limit       int       `json:"limit"`
	Remaining   int       `json:"remaining"`
	Percent     float64   `json:"percent"`
	Exceeded    bool      `json:"exceeded"`
	ResetsAt    time.Time `json:"resetsAt"`
	Exempt      bool      `json:"exempt"`
	ActiveJobID *string   `json:"activeJobId,omitempty"`
}

var ErrNullResponse = errors.New("auth failed: empty session response")
var ErrNotImplemented = errors.New("not implemented")

func decodeResponse[T comparable](body io.Reader) (T, error) {
	d := json.NewDecoder(body)
	var out T
	var e T
	if err := d.Decode(&out); err != nil {
		return out, err
	}
	if out == e {
		return out, ErrNullResponse
	}

	return out, nil
}
