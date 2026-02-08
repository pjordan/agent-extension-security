package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return fs
}

// parseInterspersed allows positional args to appear before flags by
// reordering args for the stdlib flag parser.
func parseInterspersed(fs *flag.FlagSet, args []string) error {
	return fs.Parse(reorderInterspersed(fs, args))
}

func reorderInterspersed(fs *flag.FlagSet, args []string) []string {
	reordered := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positionals = append(positionals, arg)
			continue
		}

		reordered = append(reordered, arg)

		name, hasValueInline := parseFlagToken(arg)
		if name == "" || hasValueInline || isBoolFlag(fs, name) {
			continue
		}
		if i+1 < len(args) {
			next := args[i+1]
			if next != "--" && (!strings.HasPrefix(next, "-") || next == "-") {
				reordered = append(reordered, next)
				i++
			}
		}
	}

	return append(reordered, positionals...)
}

func parseFlagToken(arg string) (name string, hasValueInline bool) {
	if strings.HasPrefix(arg, "--") {
		name = strings.TrimPrefix(arg, "--")
	} else if strings.HasPrefix(arg, "-") {
		name = strings.TrimPrefix(arg, "-")
	}
	if name == "" {
		return "", false
	}
	if idx := strings.Index(name, "="); idx >= 0 {
		return name[:idx], true
	}
	return name, false
}

func isBoolFlag(fs *flag.FlagSet, name string) bool {
	f := fs.Lookup(name)
	if f == nil {
		return false
	}
	if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); ok {
		return bf.IsBoolFlag()
	}
	return false
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
