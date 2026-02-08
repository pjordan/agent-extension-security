package main

import (
	"fmt"

	"github.com/pjordan/agent-extension-security/internal/crypto"
	"github.com/pjordan/agent-extension-security/internal/util"
)

func runSign(args []string) {
	fs := newFlagSet("sign")
	keyPath := fs.String("key", "", "dev key file (json with private key)")
	out := fs.String("out", "", "output signature file (json)")
	dieIf(fs.Parse(args))
	if fs.NArg() < 1 || *keyPath == "" || *out == "" {
		dieIf(fmt.Errorf("usage: agentsec sign <artifact.aext> --key <devkey.json> --out <sig.json>"))
	}
	art := fs.Arg(0)

	sha, err := util.Sha256File(art)
	dieIf(err)
	digest := "sha256:" + sha

	k, err := crypto.LoadDevKey(*keyPath)
	dieIf(err)
	priv, err := k.PrivateKey()
	dieIf(err)
	pub, err := k.PublicKey()
	dieIf(err)

	sig, err := crypto.SignDigest(digest, priv, pub)
	dieIf(err)
	dieIf(crypto.SaveSignature(*out, sig))
	fmt.Println("wrote signature:", *out)
}

func runVerify(args []string) {
	fs := newFlagSet("verify")
	sigPath := fs.String("sig", "", "signature file (json)")
	pubPath := fs.String("pub", "", "optional pubkey file (json) to override signature public_key")
	dieIf(fs.Parse(args))
	if fs.NArg() < 1 || *sigPath == "" {
		dieIf(fmt.Errorf("usage: agentsec verify <artifact.aext> --sig <sig.json> [--pub <pubkey.json>]"))
	}
	art := fs.Arg(0)

	sha, err := util.Sha256File(art)
	dieIf(err)
	digest := "sha256:" + sha

	sig, err := crypto.LoadSignature(*sigPath)
	dieIf(err)

	var pub []byte
	if *pubPath != "" {
		k, err := crypto.LoadDevKey(*pubPath)
		dieIf(err)
		p, err := k.PublicKey()
		dieIf(err)
		pub = p
	} else {
		// trust the public key embedded in signature (dev mode)
		k := &crypto.DevKeyFile{Type: "ed25519", Public: sig.PublicKey}
		p, err := k.PublicKey()
		dieIf(err)
		pub = p
	}

	dieIf(crypto.Verify(sig, digest, pub))
	fmt.Println("verified:", art)
}
