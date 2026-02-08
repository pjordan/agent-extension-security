package main

import (
	"fmt"
	"os"

	"github.com/pjordan/agent-extension-security/internal/crypto"
)

func runKeygen(args []string) {
	fs := newFlagSet("keygen")
	out := fs.String("out", "", "output key file (json)")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: agentsec keygen --out <file>

Generate an Ed25519 dev signing keypair. The output file contains both the
private and public keys in JSON format. Keep this file private.

Flags:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Example:
  agentsec keygen --out devkey.json
`)
	}
	dieIf(parseInterspersed(fs, args))
	if *out == "" {
		dieIf(fmt.Errorf("--out is required"))
	}
	k, err := crypto.GenerateDevKeypair()
	dieIf(err)
	dieIf(crypto.SaveKey(*out, k))
	fmt.Println("wrote key:", *out)
}
