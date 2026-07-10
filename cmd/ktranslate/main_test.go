package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootUsageIncludesDetect(t *testing.T) {
	flags := rootFlagSet()
	var out bytes.Buffer
	flags.SetOutput(&out)

	flags.Usage()

	if !strings.Contains(out.String(), "detect       detect source language") {
		t.Fatalf("usage does not include detect command:\n%s", out.String())
	}
}

func TestSplitCSV(t *testing.T) {
	got := splitCSV("eu, es, ,en")
	want := []string{"eu", "es", "en"}
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	}
}
