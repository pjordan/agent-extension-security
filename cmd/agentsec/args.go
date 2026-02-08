package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return fs
}

func dieIf(err error) {
	if err == nil {
		return
	}
	// flag package returns this on -h; treat as normal exit
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
