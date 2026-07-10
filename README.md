# go-kagi-translate

Unofficial Go client and CLI for `translate.kagi.com`.

This project uses Kagi Translate's private web API. It requires a valid Kagi
subscription and can break when Kagi changes the web application or request
contract.

A full example in the form of a cli can be found in
[cmd/ktranslate/main.go](./cmd/ktranslate/main.go).

Features:

- translate text programmatically via `gokagitranslate.(*Kagi).Translate(ctx, from, to, text string)`
- detect the language of a text with `gokagitranslate.(*Kagi).Detect(ctx, text string)`
- inspect quota usage programmatically via `gokagitranslate.(*Kagi).Quota(ctx)`
- configurable API client via `gokagitranslate.New(token).{WithClient,WithUserAgent}`
- uncomplicated opinionated oneshot translation via `gokagitranslate.OneShot(ctx, token, from, to, text string)`
- netiquette adhering with both project name and contact email in user agent, dont @ me kagi team :)
- cli tool for experimenting with all features the api provides (translation, quotas, language detection)

## Install

```sh
go get github.com/xnacly/go-kagi-translate
```

For the CLI:

```sh
go install github.com/xnacly/go-kagi-translate/cmd/ktranslate@latest
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	gokagitranslate "github.com/xnacly/go-kagi-translate"
)

func main() {
	client := gokagitranslate.New(os.Getenv("KAGI_TOKEN"))

	translation, err := client.Translate(context.Background(), "es", "en", "me llamo teo")
	if err != nil {
		panic(err)
	}
	fmt.Println(translation.Translation)

	detected, err := client.Detect(context.Background(), "agur, zer moduz zaude?")
	if err != nil {
		panic(err)
	}
	fmt.Println(detected.DetectedLanguage.Iso)

	quota, err := client.Quota(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(quota.Translate.Remaining)
}
```

For lower-level translation options, use `TranslateWithParams`:

```go
res, err := client.TranslateWithParams(context.Background(), gokagitranslate.TranslateParams{
	Text:               "hello",
	From:               "en",
	To:                 "es",
	Formality:          "default",
	SpeakerGender:      "unknown",
	AddresseeGender:    "unknown",
	LanguageComplexity: "standard",
	TranslationStyle:   "natural",
	Model:              "standard",
})
```

Language detection is best-effort, especially for short or ambiguous text. To
provide Kagi with UI-like recent language context, use `DetectWithParams`:

```go
res, err := client.DetectWithParams(context.Background(), gokagitranslate.DetectParams{
	Text:                "agur",
	IncludeAlternatives: true,
	RecentLanguages:     []string{"eu", "es", "en"},
})
```

## CLI

The CLI expects `KAGI_TOKEN` in the environment:

```sh
KAGI_TOKEN=... ktranslate translate -from es -to en "me llamo matteo"
KAGI_TOKEN=... ktranslate detect -json "agur, zer moduz zaude?"
KAGI_TOKEN=... ktranslate detect -json -recent eu,es,en "agur"
KAGI_TOKEN=... ktranslate quota
```

## Reason

I use kagi translate daily in a large fashion due to working with spanish,
american, russian and german colleagues. Since I currently live in the basque
country and kagi has high quality spanish and basque source and target
translations I wanted to use said product programmatically via a rest api, only
to discover: There is no documentation and the seemingly available beta program
is gated behind writing the support with my usecase and then I would, somehow
get an api key. Anyway, this is the reversed web api, if you have a
subscription, thus a session link and token, have fun.

## Authentication

> This requires a subscription to kagi, since translate is only available to this.

Both the example mentioned before and the tests expect the kagi private session
token in the `KAGI_TOKEN` env variable, this value can be get from the starred
(`*`) section of the session link one can request from kagi when going to the
sidebar on kagi.com and clicking on `Copy` in the `Session Link` section:

```text
https://kagi.com/search?token=***************************************************************************************&q=%s
```
