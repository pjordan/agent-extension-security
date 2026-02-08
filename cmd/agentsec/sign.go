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
	dieIf(parseInterspersed(fs, args))
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
	pubPath := fs.String("pub", "", "trusted pubkey file (json)")
	allowEmbeddedKey := fs.Bool("allow-embedded-key", false, "allow trusting signature public_key (insecure/dev-only)")
	dieIf(parseInterspersed(fs, args))
	if fs.NArg() < 1 || *sigPath == "" {
		dieIf(fmt.Errorf("usage: agentsec verify <artifact.aext> --sig <sig.json> (--pub <pubkey.json> | --allow-embedded-key)"))
	}
	art := fs.Arg(0)

	dieIf(verifyArtifact(art, *sigPath, *pubPath, *allowEmbeddedKey))
	fmt.Println("verified:", art)
}

func verifyArtifact(art, sigPath, pubPath string, allowEmbeddedKey bool) error {
	sha, err := util.Sha256File(art)
	if err != nil {
		return err
	}
	digest := "sha256:" + sha

	sig, err := crypto.LoadSignature(sigPath)
	if err != nil {
		return err
	}

	var pub []byte
	if pubPath != "" {
		k, err := crypto.LoadDevKey(pubPath)
		if err != nil {
			return err
		}
		p, err := k.PublicKey()
		if err != nil {
			return err
		}
		pub = p
	} else if allowEmbeddedKey {
		// Insecure compatibility path for local/dev workflows.
		k := &crypto.DevKeyFile{Type: "ed25519", Public: sig.PublicKey}
		p, err := k.PublicKey()
		if err != nil {
			return err
		}
		pub = p
	} else {
		return fmt.Errorf("no trusted key provided: pass --pub <pubkey.json> (or --allow-embedded-key for insecure/dev mode)")
	}

	return crypto.Verify(sig, digest, pub)
}
