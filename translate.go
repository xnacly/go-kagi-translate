package gokagitranslate

import "context"

// TranslateParams contains the request body sent to Kagi's translate endpoint.
type TranslateParams struct {
	Text                   string `json:"text"`
	From                   string `json:"from"`
	To                     string `json:"to"`
	Stream                 bool   `json:"stream"`
	PredictedLanguage      string `json:"predicted_language,omitempty"`
	Formality              string `json:"formality"`
	SpeakerGender          string `json:"speaker_gender"`
	AddresseeGender        string `json:"addressee_gender"`
	LanguageComplexity     string `json:"language_complexity"`
	TranslationStyle       string `json:"translation_style"`
	Context                string `json:"context"`
	Model                  string `json:"model"`
	SessionToken           string `json:"session_token"`
	DictionaryLanguage     string `json:"dictionary_language"`
	UseDefinitionContext   bool   `json:"use_definition_context"`
	EnableLanguageFeatures bool   `json:"enable_language_features"`
}

// TranslateWithParams translates text using the provided low-level options.
func (kt *Kagi) TranslateWithParams(ctx context.Context, params TranslateParams) (TranslateResponse, error) {
	if err := kt.auth(ctx); err != nil {
		return TranslateResponse{}, err
	}
	if params.SessionToken == "" {
		params.SessionToken = kt.session.Token
	}

	return decodeJSON[TranslateResponse](ctx, kt, "POST", translate, "translate", params)
}

// Translate translates text from one language to another using standard options.
func (kt *Kagi) Translate(ctx context.Context, from, to, text string) (TranslateResponse, error) {
	params := TranslateParams{
		Text:                   text,
		From:                   from,
		To:                     to,
		Stream:                 false,
		Formality:              "default",
		SpeakerGender:          "unknown",
		AddresseeGender:        "unknown",
		LanguageComplexity:     "standard",
		TranslationStyle:       "natural",
		Context:                "",
		Model:                  "standard",
		SessionToken:           kt.session.Token,
		DictionaryLanguage:     to,
		UseDefinitionContext:   false,
		EnableLanguageFeatures: false,
	}

	return kt.TranslateWithParams(ctx, params)
}
