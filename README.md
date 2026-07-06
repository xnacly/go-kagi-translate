# go-kagi-translate

> WARNING: this is a WIP

> This is an unofficial project using the web api translate.kagi.com, it
> requires a valid kagi subscription

A full example can be found in [cmd/main.go](./cmd/main.go). 

Features:

- translate text programmatically via `gokagitranslate.(*Kagi).Translate(ctx, from, to, text string)`
- configurable API client via `gokagitranslate.New().{WithClient,WithToken}`
- uncomplicated opinionated oneshot translation via `gokagitranslate.OneShot(ctx, token, from, to, text string)`
- netiquette adhering with both project name and contact email in user agent, dont @ me kagi team :)

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
