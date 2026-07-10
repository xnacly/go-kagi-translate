package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	gokagitranslate "github.com/xnacly/go-kagi-translate"
)

func die(err error) {
	slog.Error("failed", "err", err)
	os.Exit(1)
}

func rootFlagSet() *flag.FlagSet {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), `Usage:
  %s <command> [flags] [args]

Commands:
  translate    translate text
  detect       detect source language
  quota        show translate quota usage

Run "%s <command> -h" for command flags.
`, os.Args[0], os.Args[0])
	}
	return flags
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	root := rootFlagSet()
	if len(os.Args) == 1 {
		root.Usage()
		os.Exit(2)
	}

	args := os.Args[1:]
	root.Parse(args)
	commandArgs := root.Args()
	if len(commandArgs) == 0 {
		root.Usage()
		os.Exit(2)
	}

	switch commandArgs[0] {
	case "translate":
		runTranslate(ctx, commandArgs[1:])
	case "detect":
		runDetect(ctx, commandArgs[1:])
	case "quota":
		runQuota(ctx, commandArgs[1:])
	case "help":
		root.Usage()
	default:
		root.Usage()
		os.Exit(2)
	}
}

func runTranslate(ctx context.Context, args []string) {
	var (
		from    = ""
		to      = ""
		verbose = false
		asJSON  = false
	)
	flags := flag.NewFlagSet("translate", flag.ExitOnError)
	flags.StringVar(&from, "from", "", "set source language")
	flags.StringVar(&to, "to", "", "set target language")
	flags.BoolVar(&asJSON, "json", false, "print full response as JSON")
	flags.BoolVar(&verbose, "v", false, "verbose logging")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), `Usage:
  %s translate -from <lang> -to <lang> [flags] <text...>

Flags:
`, os.Args[0])
		flags.PrintDefaults()
	}
	flags.Parse(args)

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
	} else if len(flags.Args()) == 0 {
		die(errors.New("nothing to translate"))
	}

	client, err := newClient()
	if err != nil {
		die(err)
	}
	slog.Debug("created kagi translate client")
	output, err := client.Translate(ctx, from, to, strings.Join(flags.Args(), " "))
	if err != nil {
		die(err)
	}
	slog.Debug("translated", "from", from, "to", to)
	if asJSON {
		printJSON(output)
		return
	}
	fmt.Println(output.Translation)
}

func runDetect(ctx context.Context, args []string) {
	var (
		verbose         = false
		asJSON          = false
		recentLanguages = ""
	)
	flags := flag.NewFlagSet("detect", flag.ExitOnError)
	flags.StringVar(&recentLanguages, "recent", "", "comma-separated recent language codes")
	flags.BoolVar(&asJSON, "json", false, "print full response as JSON")
	flags.BoolVar(&verbose, "v", false, "verbose logging")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), `Usage:
  %s detect [flags] <text...>

Flags:
`, os.Args[0])
		flags.PrintDefaults()
	}
	flags.Parse(args)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if len(flags.Args()) == 0 {
		die(errors.New("nothing to detect"))
	}

	client, err := newClient()
	if err != nil {
		die(err)
	}
	slog.Debug("created kagi translate client")
	params := gokagitranslate.DetectParams{
		Text:                strings.Join(flags.Args(), " "),
		IncludeAlternatives: true,
		RecentLanguages:     splitCSV(recentLanguages),
	}
	output, err := client.DetectWithParams(ctx, params)
	if err != nil {
		die(err)
	}
	if asJSON {
		printJSON(output)
		return
	}
	fmt.Println(output.DetectedLanguage.Iso)
}

func runQuota(ctx context.Context, args []string) {
	var (
		verbose = false
		asJSON  = false
	)
	flags := flag.NewFlagSet("quota", flag.ExitOnError)
	flags.BoolVar(&asJSON, "json", false, "print full response as JSON")
	flags.BoolVar(&verbose, "v", false, "verbose logging")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), `Usage:
  %s quota [flags]

Flags:
`, os.Args[0])
		flags.PrintDefaults()
	}
	flags.Parse(args)
	if len(flags.Args()) != 0 {
		die(fmt.Errorf("quota does not accept arguments: %s", strings.Join(flags.Args(), " ")))
	}

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	client, err := newClient()
	if err != nil {
		die(err)
	}
	slog.Debug("created kagi translate client")
	quota, err := client.Quota(ctx)
	if err != nil {
		die(err)
	}
	if asJSON {
		printJSON(quota)
		return
	}
	printQuota(quota)
}

func newClient() (*gokagitranslate.Kagi, error) {
	token, ok := os.LookupEnv("KAGI_TOKEN")
	if !ok {
		return nil, errors.New("no KAGI_TOKEN env variable set")
	}
	slog.Debug("found KAGI_TOKEN env variable")
	return gokagitranslate.New(token).WithClient(&http.Client{}), nil
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		die(err)
	}
}

func printQuota(q gokagitranslate.QuotaResponse) {
	printQuotaLine(q.Translate)
	printQuotaLine(q.Proofread)
	printQuotaLine(q.Document)
}

func printQuotaLine(q gokagitranslate.Quota) {
	reset := "unknown"
	if !q.ResetsAt.IsZero() {
		reset = q.ResetsAt.Format(time.RFC3339)
	}
	exempt := ""
	if q.Exempt {
		exempt = ", exempt"
	}
	activeJob := ""
	if q.ActiveJobID != nil {
		activeJob = fmt.Sprintf(", active job %s", *q.ActiveJobID)
	}
	fmt.Printf("%s: %d/%d used, %d remaining (%.2f%%), resets %s%s%s\n",
		q.Kind, q.Used, q.Limit, q.Remaining, q.Percent, reset, exempt, activeJob)
}
