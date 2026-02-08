package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "version":
		fmt.Println(versionString())
	case "init":
		runInit(args)
	case "keygen":
		runKeygen(args)
	case "package":
		runPackage(args)
	case "manifest":
		runManifest(args)
	case "sbom":
		runSBOM(args)
	case "provenance":
		runProvenance(args)
	case "scan":
		runScan(args)
	case "sign":
		runSign(args)
	case "verify":
		runVerify(args)
	case "install":
		runInstall(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func versionString() string {
	base := fmt.Sprintf("agentsec %s", version)
	var metadata []string
	if commit != "" && commit != "none" {
		metadata = append(metadata, "commit="+commit)
	}
	if date != "" && date != "unknown" {
		metadata = append(metadata, "built="+date)
	}
	if len(metadata) == 0 {
		return base
	}
	return fmt.Sprintf("%s (%s)", base, strings.Join(metadata, ", "))
}

func usage() {
	fmt.Print(`agentsec - agent extension security (scaffold)

Usage:
  agentsec <command> [flags]

Setup:
  init                Scaffold a new extension project (manifest + key + policy)
  keygen              Generate an Ed25519 dev signing keypair

Build:
  package             Package a directory into an .aext artifact
  manifest init       Create an AEM manifest with least-privilege defaults
  manifest validate   Validate an existing AEM manifest

Attest:
  sbom                Generate a reference SBOM for an artifact
  provenance          Generate reference provenance for an artifact
  scan                Heuristic scan of an artifact for risky patterns

Security:
  sign                Sign an artifact with an Ed25519 dev key
  verify              Verify an artifact's signature

Deploy:
  install             Verify, evaluate policy, and extract an artifact

Other:
  version             Print version information

Run 'agentsec <command> -h' for help on a specific command.

Notes:
  - This initial scaffold uses local ed25519 dev keys.
  - Flags may be provided before or after positional arguments.
  - See docs/sigstore.md for how Sigstore/Cosign fits in for production signing.
`)
}
