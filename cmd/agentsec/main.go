package main

import (
	"fmt"
	"os"
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
		fmt.Println("agentsec v0.1.0 (scaffold)")
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

func usage() {
	fmt.Print(`agentsec - agent extension security (scaffold)

Usage:
  agentsec version
  agentsec keygen --out <file>
  agentsec package <dir> --out <artifact.aext>
  agentsec manifest init <dir> --id <id> --type <skill|mcp-server|plugin> --version <ver> --out <aem.json>
  agentsec manifest validate <aem.json>
  agentsec sbom <artifact.aext> --out <sbom.spdx.json>
  agentsec provenance <artifact.aext> --source-repo <url> --source-rev <rev> --out <provenance.json>
  agentsec scan <artifact.aext> --out <scan.json>
  agentsec sign <artifact.aext> --key <devkey.json> --out <sig.json>
  agentsec verify <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key)
  agentsec install <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key) --aem <aem.json> --policy <policy.json> --dest <dir>

Notes:
  - This initial scaffold uses local ed25519 dev keys.
  - See docs/sigstore.md for how Sigstore/Cosign fits in for production signing.
`)
}
