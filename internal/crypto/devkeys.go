package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type DevKeyFile struct {
	Type    string `json:"type"`
	Public  string `json:"public"`            // base64
	Private string `json:"private,omitempty"` // base64
}

func GenerateDevKeypair() (*DevKeyFile, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &DevKeyFile{
		Type:    "ed25519",
		Public:  base64.StdEncoding.EncodeToString(pub),
		Private: base64.StdEncoding.EncodeToString(priv),
	}, nil
}

func LoadDevKey(path string) (*DevKeyFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var k DevKeyFile
	if err := json.Unmarshal(b, &k); err != nil {
		return nil, err
	}
	if k.Type != "ed25519" {
		return nil, fmt.Errorf("unsupported key type: %s", k.Type)
	}
	if k.Public == "" {
		return nil, fmt.Errorf("missing public key")
	}
	return &k, nil
}

func (k *DevKeyFile) PublicKey() (ed25519.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(k.Public)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(b), nil
}

func (k *DevKeyFile) PrivateKey() (ed25519.PrivateKey, error) {
	if k.Private == "" {
		return nil, fmt.Errorf("no private key in key file")
	}
	b, err := base64.StdEncoding.DecodeString(k.Private)
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(b), nil
}

func SaveKey(path string, k *DevKeyFile) error {
	b, err := json.MarshalIndent(k, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o600)
}
