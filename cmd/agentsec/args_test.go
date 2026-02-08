package main

import "testing"

func TestParseInterspersedAllowsPositionalBeforeFlags(t *testing.T) {
	fs := newFlagSet("package")
	out := fs.String("out", "", "")

	err := parseInterspersed(fs, []string{"./examples/skills/hello-world", "--out", "./_demo/hello-world.aext"})
	if err != nil {
		t.Fatalf("parseInterspersed() error = %v", err)
	}
	if got, want := fs.NArg(), 1; got != want {
		t.Fatalf("NArg = %d, want %d", got, want)
	}
	if got, want := fs.Arg(0), "./examples/skills/hello-world"; got != want {
		t.Fatalf("Arg(0) = %q, want %q", got, want)
	}
	if got, want := *out, "./_demo/hello-world.aext"; got != want {
		t.Fatalf("out = %q, want %q", got, want)
	}
}

func TestParseInterspersedKeepsBoolFlagsWithoutValue(t *testing.T) {
	fs := newFlagSet("verify")
	sig := fs.String("sig", "", "")
	allowEmbedded := fs.Bool("allow-embedded-key", false, "")

	err := parseInterspersed(fs, []string{"artifact.aext", "--sig", "sig.json", "--allow-embedded-key"})
	if err != nil {
		t.Fatalf("parseInterspersed() error = %v", err)
	}
	if got, want := fs.NArg(), 1; got != want {
		t.Fatalf("NArg = %d, want %d", got, want)
	}
	if got, want := fs.Arg(0), "artifact.aext"; got != want {
		t.Fatalf("Arg(0) = %q, want %q", got, want)
	}
	if got, want := *sig, "sig.json"; got != want {
		t.Fatalf("sig = %q, want %q", got, want)
	}
	if !*allowEmbedded {
		t.Fatal("allow-embedded-key = false, want true")
	}
}
