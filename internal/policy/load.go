package policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Load(path string) (*Policy, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load policy %s: %w", path, err)
	}
	var p Policy
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("load policy %s: %w", path, err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("load policy %s: policy must contain a single JSON object", path)
	}
	if p.Mode == "" {
		p.Mode = "enforce"
	}
	if p.Mode != "enforce" && p.Mode != "warn" {
		return nil, fmt.Errorf("load policy %s: mode must be enforce|warn", path)
	}
	return &p, nil
}
