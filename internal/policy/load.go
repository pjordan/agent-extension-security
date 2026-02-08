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
		return nil, err
	}
	var p Policy
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return nil, err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("policy must contain a single JSON object")
	}
	if p.Mode == "" {
		p.Mode = "enforce"
	}
	if p.Mode != "enforce" && p.Mode != "warn" {
		return nil, fmt.Errorf("mode must be enforce|warn")
	}
	return &p, nil
}
