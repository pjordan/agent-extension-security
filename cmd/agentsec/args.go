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

func requireArg(args []string, i int, label string) (string, error) {
    if len(args) <= i || args[i] == "" {
        return "", fmt.Errorf("missing %s", label)
    }
    return args[i], nil
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
