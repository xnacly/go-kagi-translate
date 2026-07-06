package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gokagitranslate "github.com/xnacly/go-kagi-translate"
)

func die(err error) {
	slog.Error("failed", "err", err)
	os.Exit(1)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
	}
	slog.Debug("parsed arguments",
		"from", from,
		"to", to)

	if len(from) == 0 {
		die(errors.New("no -from defined"))
	} else if len(to) == 0 {
		die(errors.New("no -to defined"))
	} else if len(flag.Args()) == 0 {
		die(errors.New("nothing to translate"))
	}

	token, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		die(errors.New("no KAGI_TOKEN env variable set"))
	}
	slog.Debug("found KAGI_TOKEN env variable")

	client := gokagitranslate.New().WithClient(&http.Client{}).WithToken(token)
	slog.Debug("created kagi translate client")
	output, err := client.Translate(ctx, from, to, strings.Join(flag.Args(), " "))
	if err != nil {
		die(err)
	}
	slog.Debug("translated", "from", from, "to", to)
	fmt.Println(output)
}
