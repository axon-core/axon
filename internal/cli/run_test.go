package cli

import (
	"testing"
)

func TestParseMCPServerFlags(t *testing.T) {
	tests := []struct {
		name    string
		flags   []string
		want    map[string]struct{ typ, url, command string }
		wantErr bool
	}{
		{
			name:  "Empty flags",
			flags: nil,
			want:  nil,
		},
		{
			name:  "HTTP server",
			flags: []string{"my-api=http:https://api.example.com/mcp"},
			want: map[string]struct{ typ, url, command string }{
				"my-api": {typ: "http", url: "https://api.example.com/mcp"},
			},
		},
		{
			name:  "SSE server",
			flags: []string{"my-sse=sse:https://api.example.com/sse"},
			want: map[string]struct{ typ, url, command string }{
				"my-sse": {typ: "sse", url: "https://api.example.com/sse"},
			},
		},
		{
			name:  "Stdio server",
			flags: []string{"local-tool=stdio:npx -y @example/server"},
			want: map[string]struct{ typ, url, command string }{
				"local-tool": {typ: "stdio", command: "npx"},
			},
		},
		{
			name: "Multiple servers",
			flags: []string{
				"api=http:https://api.example.com/mcp",
				"tool=stdio:python server.py",
			},
			want: map[string]struct{ typ, url, command string }{
				"api":  {typ: "http", url: "https://api.example.com/mcp"},
				"tool": {typ: "stdio", command: "python"},
			},
		},
		{
			name:    "Missing equals sign",
			flags:   []string{"invalid-format"},
			wantErr: true,
		},
		{
			name:    "Missing colon in type:target",
			flags:   []string{"name=invalidformat"},
			wantErr: true,
		},
		{
			name:    "Invalid transport type",
			flags:   []string{"name=grpc:localhost:8080"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMCPServerFlags(tt.flags)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d servers, want %d", len(got), len(tt.want))
			}

			for name, wantServer := range tt.want {
				gotServer, ok := got[name]
				if !ok {
					t.Errorf("missing server %q", name)
					continue
				}
				if gotServer.Type != wantServer.typ {
					t.Errorf("server %q type = %q, want %q", name, gotServer.Type, wantServer.typ)
				}
				if wantServer.url != "" && gotServer.URL != wantServer.url {
					t.Errorf("server %q url = %q, want %q", name, gotServer.URL, wantServer.url)
				}
				if wantServer.command != "" && gotServer.Command != wantServer.command {
					t.Errorf("server %q command = %q, want %q", name, gotServer.Command, wantServer.command)
				}
			}
		})
	}
}

func TestParseMCPServerFlags_StdioArgs(t *testing.T) {
	flags := []string{"tool=stdio:npx -y @example/server --port 8080"}
	got, err := parseMCPServerFlags(flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	server := got["tool"]
	if server.Command != "npx" {
		t.Errorf("command = %q, want %q", server.Command, "npx")
	}
	if len(server.Args) != 4 {
		t.Fatalf("expected 4 args, got %d: %v", len(server.Args), server.Args)
	}
	expectedArgs := []string{"-y", "@example/server", "--port", "8080"}
	for i, arg := range expectedArgs {
		if server.Args[i] != arg {
			t.Errorf("arg[%d] = %q, want %q", i, server.Args[i], arg)
		}
	}
}
