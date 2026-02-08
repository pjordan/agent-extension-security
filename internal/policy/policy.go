package policy

import (
	"fmt"

	"github.com/pjordan/agent-extension-security/internal/manifest"
)

// Policy is intentionally minimal in the scaffold.
// It supports a few deny rules and can be expanded.
type Policy struct {
	Mode string `json:"mode"` // enforce|warn

	Permissions struct {
		Deny struct {
			Network struct {
				AllowIPLiterals *bool `json:"allow_ip_literals,omitempty"`
			} `json:"network,omitempty"`
			Process struct {
				AllowShell *bool `json:"allow_shell,omitempty"`
			} `json:"process,omitempty"`
		} `json:"deny,omitempty"`
	} `json:"permissions,omitempty"`
}

func Evaluate(p *Policy, aem *manifest.AEM) []string {
	var findings []string

	if p.Permissions.Deny.Network.AllowIPLiterals != nil {
		if aem.Permissions.Network.AllowIPLiterals == *p.Permissions.Deny.Network.AllowIPLiterals {
			findings = append(findings, fmt.Sprintf("denied: network.allow_ip_literals=%v", aem.Permissions.Network.AllowIPLiterals))
		}
	}
	if p.Permissions.Deny.Process.AllowShell != nil {
		if aem.Permissions.Process.AllowShell == *p.Permissions.Deny.Process.AllowShell {
			findings = append(findings, fmt.Sprintf("denied: process.allow_shell=%v", aem.Permissions.Process.AllowShell))
		}
	}
	return findings
}

// DefaultPermissivePolicy returns a warn-only policy with no deny rules.
// Used by --dev mode to allow installation without a policy file.
func DefaultPermissivePolicy() *Policy {
	return &Policy{Mode: "warn"}
}
