package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// AEM is an initial, intentionally small Agent Extension Manifest.
// This scaffold uses JSON (std library only) to keep the repo easy to build.
type AEM struct {
	Schema     string `json:"schema"`
	ID         string `json:"id"`
	Type       string `json:"type"` // skill | mcp-server | plugin
	Version    string `json:"version"`
	SourceRepo string `json:"source_repo,omitempty"`
	SourceRev  string `json:"source_rev,omitempty"`

	Permissions Permissions `json:"permissions"`
}

type Permissions struct {
	Files   FilePerms   `json:"files,omitempty"`
	Network NetPerms    `json:"network,omitempty"`
	Process ProcessPerm `json:"process,omitempty"`
}

type FilePerms struct {
	Read  []string `json:"read,omitempty"`
	Write []string `json:"write,omitempty"`
}

type NetPerms struct {
	Domains         []string `json:"domains,omitempty"`
	AllowIPLiterals bool     `json:"allow_ip_literals,omitempty"`
}

type ProcessPerm struct {
	AllowShell      bool `json:"allow_shell,omitempty"`
	AllowSubprocess bool `json:"allow_subprocess,omitempty"`
}

var semverPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)

func LoadAEM(path string) (*AEM, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load manifest %s: %w", path, err)
	}
	var m AEM
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("load manifest %s: %w", path, err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("load manifest %s: manifest must contain a single JSON object", path)
	}
	return &m, nil
}

func (m *AEM) Validate() error {
	if m.Schema != "aessf.dev/aem/v0" {
		return fmt.Errorf("schema must be aessf.dev/aem/v0")
	}
	if m.ID == "" {
		return fmt.Errorf("id is required")
	}
	switch m.Type {
	case "skill", "mcp-server", "plugin":
	default:
		return fmt.Errorf("type must be skill|mcp-server|plugin")
	}
	if !semverPattern.MatchString(m.Version) {
		return fmt.Errorf("version must be semver (e.g., 1.2.3)")
	}
	return nil
}

func (m *AEM) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}
