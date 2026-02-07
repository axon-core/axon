package controller

import (
	"encoding/json"
	"testing"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

func TestBuildMCPConfig(t *testing.T) {
	tests := []struct {
		name    string
		servers map[string]axonv1alpha1.MCPServer
		want    mcpConfigJSON
	}{
		{
			name: "HTTP server",
			servers: map[string]axonv1alpha1.MCPServer{
				"my-api": {
					Type: "http",
					URL:  "https://api.example.com/mcp",
				},
			},
			want: mcpConfigJSON{
				MCPServers: map[string]mcpServerJSON{
					"my-api": {
						Type: "http",
						URL:  "https://api.example.com/mcp",
					},
				},
			},
		},
		{
			name: "SSE server with headers",
			servers: map[string]axonv1alpha1.MCPServer{
				"sse-api": {
					Type: "sse",
					URL:  "https://api.example.com/sse",
					Headers: map[string]string{
						"Authorization": "Bearer token",
					},
				},
			},
			want: mcpConfigJSON{
				MCPServers: map[string]mcpServerJSON{
					"sse-api": {
						Type: "sse",
						URL:  "https://api.example.com/sse",
						Headers: map[string]string{
							"Authorization": "Bearer token",
						},
					},
				},
			},
		},
		{
			name: "Stdio server",
			servers: map[string]axonv1alpha1.MCPServer{
				"local-tool": {
					Type:    "stdio",
					Command: "npx",
					Args:    []string{"-y", "@example/server"},
				},
			},
			want: mcpConfigJSON{
				MCPServers: map[string]mcpServerJSON{
					"local-tool": {
						Type:    "stdio",
						Command: "npx",
						Args:    []string{"-y", "@example/server"},
					},
				},
			},
		},
		{
			name: "Stdio server with env",
			servers: map[string]axonv1alpha1.MCPServer{
				"db-tool": {
					Type:    "stdio",
					Command: "npx",
					Args:    []string{"-y", "@example/db-server"},
					Env: map[string]string{
						"DB_URL": "postgres://localhost/mydb",
					},
				},
			},
			want: mcpConfigJSON{
				MCPServers: map[string]mcpServerJSON{
					"db-tool": {
						Type:    "stdio",
						Command: "npx",
						Args:    []string{"-y", "@example/db-server"},
						Env: map[string]string{
							"DB_URL": "postgres://localhost/mydb",
						},
					},
				},
			},
		},
		{
			name: "Multiple servers",
			servers: map[string]axonv1alpha1.MCPServer{
				"api": {
					Type: "http",
					URL:  "https://api.example.com/mcp",
				},
				"local": {
					Type:    "stdio",
					Command: "npx",
					Args:    []string{"-y", "@example/local"},
				},
			},
			want: mcpConfigJSON{
				MCPServers: map[string]mcpServerJSON{
					"api": {
						Type: "http",
						URL:  "https://api.example.com/mcp",
					},
					"local": {
						Type:    "stdio",
						Command: "npx",
						Args:    []string{"-y", "@example/local"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildMCPConfig(tt.servers)
			if err != nil {
				t.Fatalf("buildMCPConfig() error = %v", err)
			}

			var result mcpConfigJSON
			if err := json.Unmarshal([]byte(got), &result); err != nil {
				t.Fatalf("failed to parse JSON output: %v", err)
			}

			if len(result.MCPServers) != len(tt.want.MCPServers) {
				t.Errorf("buildMCPConfig() got %d servers, want %d", len(result.MCPServers), len(tt.want.MCPServers))
			}

			for name, wantServer := range tt.want.MCPServers {
				gotServer, ok := result.MCPServers[name]
				if !ok {
					t.Errorf("buildMCPConfig() missing server %q", name)
					continue
				}
				if gotServer.Type != wantServer.Type {
					t.Errorf("server %q type = %q, want %q", name, gotServer.Type, wantServer.Type)
				}
				if gotServer.URL != wantServer.URL {
					t.Errorf("server %q url = %q, want %q", name, gotServer.URL, wantServer.URL)
				}
				if gotServer.Command != wantServer.Command {
					t.Errorf("server %q command = %q, want %q", name, gotServer.Command, wantServer.Command)
				}
			}
		})
	}
}

func TestBuildClaudeCodeJobWithMCPServers(t *testing.T) {
	builder := NewJobBuilder()

	task := &axonv1alpha1.Task{}
	task.Name = "test-task"
	task.Namespace = "default"
	task.Spec = axonv1alpha1.TaskSpec{
		Type:   AgentTypeClaudeCode,
		Prompt: "Test prompt",
		Credentials: axonv1alpha1.Credentials{
			Type:      axonv1alpha1.CredentialTypeAPIKey,
			SecretRef: axonv1alpha1.SecretReference{Name: "test-secret"},
		},
		MCPServers: map[string]axonv1alpha1.MCPServer{
			"my-api": {
				Type: "http",
				URL:  "https://api.example.com/mcp",
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify --mcp-config is in the args
	container := job.Spec.Template.Spec.Containers[0]
	foundMCPConfig := false
	for i, arg := range container.Args {
		if arg == "--mcp-config" {
			foundMCPConfig = true
			if i+1 >= len(container.Args) {
				t.Fatal("--mcp-config flag has no value")
			}
			expected := MCPConfigMountPath + "/" + MCPConfigFileName
			if container.Args[i+1] != expected {
				t.Errorf("--mcp-config value = %q, want %q", container.Args[i+1], expected)
			}
			break
		}
	}
	if !foundMCPConfig {
		t.Error("--mcp-config flag not found in container args")
	}

	// Verify init container
	if len(job.Spec.Template.Spec.InitContainers) != 1 {
		t.Fatalf("expected 1 init container, got %d", len(job.Spec.Template.Spec.InitContainers))
	}
	initContainer := job.Spec.Template.Spec.InitContainers[0]
	if initContainer.Name != "mcp-config" {
		t.Errorf("init container name = %q, want %q", initContainer.Name, "mcp-config")
	}

	// Verify volume
	if len(job.Spec.Template.Spec.Volumes) != 1 {
		t.Fatalf("expected 1 volume, got %d", len(job.Spec.Template.Spec.Volumes))
	}
	if job.Spec.Template.Spec.Volumes[0].Name != MCPConfigVolumeName {
		t.Errorf("volume name = %q, want %q", job.Spec.Template.Spec.Volumes[0].Name, MCPConfigVolumeName)
	}

	// Verify volume mount on main container (read-only)
	if len(container.VolumeMounts) != 1 {
		t.Fatalf("expected 1 volume mount, got %d", len(container.VolumeMounts))
	}
	if !container.VolumeMounts[0].ReadOnly {
		t.Error("MCP config volume mount should be read-only")
	}
}

func TestBuildClaudeCodeJobWithoutMCPServers(t *testing.T) {
	builder := NewJobBuilder()

	task := &axonv1alpha1.Task{}
	task.Name = "test-task"
	task.Namespace = "default"
	task.Spec = axonv1alpha1.TaskSpec{
		Type:   AgentTypeClaudeCode,
		Prompt: "Test prompt",
		Credentials: axonv1alpha1.Credentials{
			Type:      axonv1alpha1.CredentialTypeAPIKey,
			SecretRef: axonv1alpha1.SecretReference{Name: "test-secret"},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	// Verify no --mcp-config in args
	container := job.Spec.Template.Spec.Containers[0]
	for _, arg := range container.Args {
		if arg == "--mcp-config" {
			t.Error("--mcp-config flag should not be present when no MCP servers are configured")
		}
	}

	// Verify no init containers
	if len(job.Spec.Template.Spec.InitContainers) != 0 {
		t.Errorf("expected 0 init containers, got %d", len(job.Spec.Template.Spec.InitContainers))
	}

	// Verify no volumes
	if len(job.Spec.Template.Spec.Volumes) != 0 {
		t.Errorf("expected 0 volumes, got %d", len(job.Spec.Template.Spec.Volumes))
	}
}
