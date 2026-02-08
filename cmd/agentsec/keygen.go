package main

import (
	"fmt"

	"github.com/pjordan/agent-extension-security/internal/crypto"
)

func runKeygen(args []string) {
	fs := newFlagSet("keygen")
	out := fs.String("out", "", "output key file (json)")
	dieIf(fs.Parse(args))
	if *out == "" {
		dieIf(fmt.Errorf("--out is required"))
	}
	k, err := crypto.GenerateDevKeypair()
	dieIf(err)
	dieIf(crypto.SaveKey(*out, k))
	fmt.Println("wrote key:", *out)
}
