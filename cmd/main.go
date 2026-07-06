package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	gokagitranslate "github.com/xnacly/go-kagi-translate"
)

func die(err error) {
	slog.Error("kagi-translate", "err", err)
	panic(err.Error())
}

func main() {
	var (
		from    = ""
		to      = ""
		verbose = false
	)
	flag.StringVar(&from, "from", "", "set source language")
	flag.StringVar(&to, "to", "", "set target language")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("kagi-translate:debug", "msg", "parsed arguments",
			"from", from,
			"to", to,
			"verbose", verbose)
	}

	if len(from) == 0 {
		die(errors.New("no -from defined"))
	} else if len(to) == 0 {
		die(errors.New("no -to defined"))
	} else if len(flag.Args()) == 0 {
		die(errors.New("nothing to translate"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	token, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		die(errors.New("no KAGI_TOKEN env variable set"))
	}
	if verbose {
		slog.Debug("kagi-translate:debug", "msg", "found KAGI_TOKEN env variabel")
	}

	client := gokagitranslate.New().WithClient(http.Client{}).WithCtx(ctx).WithToken(token)
	slog.Debug("kagi-translate:debug", "msg", "created kagi translate client")
	if _, err := client.Auth(); err != nil {
		die(err)
	}
	slog.Debug("kagi-translate:debug", "msg", "authenticated to kagi translate")

	output, err := client.Translate(from, to, strings.Join(flag.Args(), " "))
	if err != nil {
		die(err)
	}
	slog.Debug("kagi-translate:debug", "msg", "translated", "from", from, "to", to)
	fmt.Println(output)
}
