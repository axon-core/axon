package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMCPFlag(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantName   string
		wantType   string
		wantURL    string
		wantCmd    string
		wantErr    bool
		wantErrStr string
	}{
		{
			name:     "http server",
			input:    `github={"type":"http","url":"https://api.githubcopilot.com/mcp/"}`,
			wantName: "github",
			wantType: "http",
			wantURL:  "https://api.githubcopilot.com/mcp/",
		},
		{
			name:     "stdio server",
			input:    `db={"type":"stdio","command":"npx","args":["-y","@bytebase/dbhub"]}`,
			wantName: "db",
			wantType: "stdio",
			wantCmd:  "npx",
		},
		{
			name:     "sse server",
			input:    `asana={"type":"sse","url":"https://mcp.asana.com/sse"}`,
			wantName: "asana",
			wantType: "sse",
			wantURL:  "https://mcp.asana.com/sse",
		},
		{
			name:       "missing type field",
			input:      `test={"url":"https://example.com"}`,
			wantErr:    true,
			wantErrStr: `"type" field is required`,
		},
		{
			name:       "invalid JSON",
			input:      `test=not-json`,
			wantErr:    true,
			wantErrStr: "invalid --mcp",
		},
		{
			name:       "missing name",
			input:      `={"type":"http"}`,
			wantErr:    true,
			wantErrStr: "invalid --mcp value",
		},
		{
			name:       "no equals sign",
			input:      `just-a-name`,
			wantErr:    true,
			wantErrStr: "invalid --mcp value",
		},
		{
			name:       "empty string",
			input:      ``,
			wantErr:    true,
			wantErrStr: "invalid --mcp value",
		},
		{
			name:       "unsupported type",
			input:      `test={"type":"websocket","url":"https://example.com"}`,
			wantErr:    true,
			wantErrStr: "unsupported type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := parseMCPFlag(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.wantErrStr)
				}
				if tt.wantErrStr != "" && !strings.Contains(err.Error(), tt.wantErrStr) {
					t.Errorf("Expected error containing %q, got: %v", tt.wantErrStr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if spec.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", spec.Name, tt.wantName)
			}
			if spec.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", spec.Type, tt.wantType)
			}
			if tt.wantURL != "" && spec.URL != tt.wantURL {
				t.Errorf("URL = %q, want %q", spec.URL, tt.wantURL)
			}
			if tt.wantCmd != "" && spec.Command != tt.wantCmd {
				t.Errorf("Command = %q, want %q", spec.Command, tt.wantCmd)
			}
		})
	}
}

func TestParseMCPFlag_FileRef(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "mcp.json")
	if err := os.WriteFile(f, []byte(`{"type":"http","url":"https://example.com/mcp"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	spec, err := parseMCPFlag("test=@" + f)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if spec.Name != "test" {
		t.Errorf("Name = %q, want %q", spec.Name, "test")
	}
	if spec.Type != "http" {
		t.Errorf("Type = %q, want %q", spec.Type, "http")
	}
	if spec.URL != "https://example.com/mcp" {
		t.Errorf("URL = %q, want %q", spec.URL, "https://example.com/mcp")
	}
}

func TestParseMCPFlag_HeadersAndEnv(t *testing.T) {
	input := `secure={"type":"http","url":"https://mcp.example.com","headers":{"Authorization":"Bearer tok"},"env":{"KEY":"val"}}`
	spec, err := parseMCPFlag(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if spec.Headers["Authorization"] != "Bearer tok" {
		t.Errorf("Headers = %v, want Authorization header", spec.Headers)
	}
	if spec.Env["KEY"] != "val" {
		t.Errorf("Env = %v, want KEY=val", spec.Env)
	}
}
