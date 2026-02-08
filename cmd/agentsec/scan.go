package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type ScanReport struct {
	Schema    string    `json:"schema"`
	CreatedAt time.Time `json:"created_at"`
	Findings  []Finding `json:"findings"`
	RiskScore int       `json:"risk_score"`
	Notes     string    `json:"notes"`
}

type Finding struct {
	Severity string `json:"severity"`
	Kind     string `json:"kind"`
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
}

func runScan(args []string) {
	fs := newFlagSet("scan")
	out := fs.String("out", "", "output scan report json path")
	dieIf(parseInterspersed(fs, args))
	if fs.NArg() < 1 || *out == "" {
		dieIf(fmt.Errorf("usage: agentsec scan <artifact.aext> --out <scan.json>"))
	}
	art := fs.Arg(0)

	zr, err := zip.OpenReader(art)
	dieIf(err)
	defer zr.Close()

	var rep ScanReport
	rep.Schema = "aessf.dev/attestation/scan/v0"
	rep.CreatedAt = time.Now().UTC()
	rep.Notes = "This is a lightweight, heuristic scan. Add real scanning in production."
	rep.Findings = []Finding{}
	rep.RiskScore = 0

	// Very simple heuristics for instruction-based risks.
	// Goal: surface risk to humans and policy engines; not "perfect detection".
	reCurlPipe := regexp.MustCompile(`(?i)curl\s+[^\n]+\|\s*(sh|bash)`)
	reWgetPipe := regexp.MustCompile(`(?i)wget\s+[^\n]+\|\s*(sh|bash)`)
	reBase64Exec := regexp.MustCompile(`(?i)base64\s+-d\s*\|\s*(sh|bash)`)
	rePowershellIEX := regexp.MustCompile(`(?i)powershell\s+.*iex\b`)
	reNpmGlobal := regexp.MustCompile(`(?i)npm\s+install\s+-g\b`)
	rePipInstall := regexp.MustCompile(`(?i)pip\s+install\b`)

	for _, f := range zr.File {
		name := f.Name
		lower := strings.ToLower(name)
		if !strings.HasSuffix(lower, "skill.md") && !strings.HasSuffix(lower, ".sh") && !strings.HasSuffix(lower, ".ps1") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			continue
		}
		b, _ := ioReadAllLimit(rc, 256*1024)
		rc.Close()
		s := string(b)

		// check patterns
		check := func(re *regexp.Regexp, sev, kind, msg string) {
			loc := re.FindStringIndex(s)
			if loc == nil {
				return
			}
			start := loc[0]
			end := loc[1]
			snippet := s[start:end]
			rep.Findings = append(rep.Findings, Finding{
				Severity: sev,
				Kind:     kind,
				Message:  msg,
				File:     name,
				Snippet:  snippet,
			})
			switch sev {
			case "high":
				rep.RiskScore += 30
			case "medium":
				rep.RiskScore += 10
			default:
				rep.RiskScore += 3
			}
		}

		check(reCurlPipe, "high", "exec_remote_script", "curl | sh pattern detected")
		check(reWgetPipe, "high", "exec_remote_script", "wget | sh pattern detected")
		check(reBase64Exec, "high", "obfuscation_exec", "base64 decode piped to shell detected")
		check(rePowershellIEX, "high", "powershell_iex", "PowerShell IEX pattern detected")
		check(reNpmGlobal, "medium", "package_install", "global npm install in instructions/scripts detected")
		check(rePipInstall, "medium", "package_install", "pip install in instructions/scripts detected")
	}

	// cap score
	if rep.RiskScore > 100 {
		rep.RiskScore = 100
	}

	outBytes, err := json.MarshalIndent(&rep, "", "  ")
	dieIf(err)
	dieIf(os.WriteFile(*out, append(outBytes, '\n'), 0o644))
	fmt.Println("wrote scan report:", *out)
}
