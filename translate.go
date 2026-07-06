package gokagitranslate

type TranslateParams struct {
	Text                     string
	From                     string
	To                       string
	Stream                   bool
	Formality                string
	Speaker_gender           string
	Addressee_gender         string
	Language_complexity      string
	Translation_style        string
	Context                  string
	Model                    string
	Session_token            string
	Dictionary_language      string
	Use_definition_context   bool
	Enable_language_features bool
}

func (kt *Kagi) TranslateWithParams(params TranslateParams) (TranslateResponse, error) {
	return TranslateResponse{}, nil
}

func (kt *Kagi) Translate(from, to, text string) (TranslateResponse, error) {
	params := TranslateParams{
		Text:                     text,
		From:                     from,
		To:                       to,
		Stream:                   false,
		Formality:                "default",
		Speaker_gender:           "unknown",
		Addressee_gender:         "unknown",
		Language_complexity:      "standard",
		Translation_style:        "natural",
		Context:                  "",
		Model:                    "standard",
		Session_token:            kt.session,
		Dictionary_language:      to,
		Use_definition_context:   false,
		Enable_language_features: false,
	}

	return kt.TranslateWithParams(params)
}
