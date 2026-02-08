package policy

import (
	"encoding/json"
	"os"
)

func Load(path string) (*Policy, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Policy
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	if p.Mode == "" {
		p.Mode = "enforce"
	}
	return &p, nil
}
