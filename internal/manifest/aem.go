package manifest

import (
    "encoding/json"
    "fmt"
    "os"
)

// AEM is an initial, intentionally small Agent Extension Manifest.
// This scaffold uses JSON (std library only) to keep the repo easy to build.
type AEM struct {
    Schema     string `json:"schema"`
    ID         string `json:"id"`
    Type       string `json:"type"`    // skill | mcp-server | plugin
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

func LoadAEM(path string) (*AEM, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var m AEM
    if err := json.Unmarshal(b, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

func (m *AEM) Validate() error {
    if m.Schema == "" {
        return fmt.Errorf("schema is required")
    }
    if m.ID == "" {
        return fmt.Errorf("id is required")
    }
    switch m.Type {
    case "skill", "mcp-server", "plugin":
    default:
        return fmt.Errorf("type must be skill|mcp-server|plugin")
    }
    if m.Version == "" {
        return fmt.Errorf("version is required")
    }
    return nil
}

func (m *AEM) ToJSON() ([]byte, error) {
    return json.MarshalIndent(m, "", "  ")
}
