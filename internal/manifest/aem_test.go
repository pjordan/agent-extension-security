package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func validAEM() *AEM {
	return &AEM{
		Schema:  "aessf.dev/aem/v0",
		ID:      "com.example.test",
		Type:    "skill",
		Version: "0.1.0",
	}
}

func TestAEMValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutator func(*AEM)
		wantErr bool
	}{
		{
			name:    "valid",
			mutator: func(_ *AEM) {},
			wantErr: false,
		},
		{
			name: "missing schema",
			mutator: func(a *AEM) {
				a.Schema = ""
			},
			wantErr: true,
		},
		{
			name: "missing id",
			mutator: func(a *AEM) {
				a.ID = ""
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			mutator: func(a *AEM) {
				a.Type = "widget"
			},
			wantErr: true,
		},
		{
			name: "missing version",
			mutator: func(a *AEM) {
				a.Version = ""
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := validAEM()
			tc.mutator(a)
			err := a.Validate()
			if tc.wantErr && err == nil {
				t.Fatal("Validate() expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestLoadAEM(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "aem.json")
		if err := os.WriteFile(p, []byte("{bad"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := LoadAEM(p); err == nil {
			t.Fatal("LoadAEM() expected error, got nil")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		a := validAEM()
		b, err := a.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON() error = %v", err)
		}

		p := filepath.Join(t.TempDir(), "aem.json")
		if err := os.WriteFile(p, b, 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		loaded, err := LoadAEM(p)
		if err != nil {
			t.Fatalf("LoadAEM() error = %v", err)
		}
		if err := loaded.Validate(); err != nil {
			t.Fatalf("loaded Validate() error = %v", err)
		}

		raw := map[string]any{}
		if err := json.Unmarshal(b, &raw); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		if raw["id"] != a.ID {
			t.Fatalf("json id = %v, want %s", raw["id"], a.ID)
		}
	})
}
