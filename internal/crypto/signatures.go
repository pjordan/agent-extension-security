    package crypto

    import (
        "crypto/ed25519"
        "encoding/base64"
        "encoding/json"
        "fmt"
        "os"
        "time"
    )

    type Signature struct {
        Alg       string `json:"alg"`        // ed25519
        Digest    string `json:"digest"`     // sha256:<hex>
        Sig       string `json:"sig"`        // base64 signature bytes
        PublicKey string `json:"public_key"` // base64 ed25519 pub
        CreatedAt string `json:"created_at"` // RFC3339
    }

    func SignDigest(digest string, priv ed25519.PrivateKey, pub ed25519.PublicKey) (*Signature, error) {
        sig := ed25519.Sign(priv, []byte(digest))
        return &Signature{
            Alg:       "ed25519",
            Digest:    digest,
            Sig:       base64.StdEncoding.EncodeToString(sig),
            PublicKey: base64.StdEncoding.EncodeToString(pub),
            CreatedAt: time.Now().UTC().Format(time.RFC3339),
        }, nil
    }

    func Verify(sig *Signature, digest string, pub ed25519.PublicKey) error {
        if sig.Alg != "ed25519" {
            return fmt.Errorf("unsupported signature alg: %s", sig.Alg)
        }
        if sig.Digest != digest {
            return fmt.Errorf("digest mismatch: sig=%s want=%s", sig.Digest, digest)
        }
        b, err := base64.StdEncoding.DecodeString(sig.Sig)
        if err != nil {
            return err
        }
        if !ed25519.Verify(pub, []byte(digest), b) {
            return fmt.Errorf("signature verification failed")
        }
        return nil
    }

    func LoadSignature(path string) (*Signature, error) {
        b, err := os.ReadFile(path)
        if err != nil {
            return nil, err
        }
        var s Signature
        if err := json.Unmarshal(b, &s); err != nil {
            return nil, err
        }
        return &s, nil
    }

    func SaveSignature(path string, s *Signature) error {
        b, err := json.MarshalIndent(s, "", "  ")
        if err != nil {
            return err
        }
        return os.WriteFile(path, append(b, '\n'), 0o644)
    }
