package controller

import (
	"testing"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

func TestBuildClaudeCodeJob(t *testing.T) {
	tests := []struct {
		name          string
		task          *axonv1alpha1.Task
		workspace     *axonv1alpha1.WorkspaceSpec
		wantImage     string
		wantCommand   []string
		wantArgs      []string
		wantEnvNames  []string
		wantEnvValues map[string]string // for non-secret env vars
	}{
		{
			name: "Default image with API key",
			task: &axonv1alpha1.Task{
				Spec: axonv1alpha1.TaskSpec{
					Type:   AgentTypeClaudeCode,
					Prompt: "Hello world",
					Credentials: axonv1alpha1.Credentials{
						Type:      axonv1alpha1.CredentialTypeAPIKey,
						SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
					},
				},
			},
			wantImage:   ClaudeCodeImage,
			wantCommand: []string{AgentEntrypoint},
			wantArgs:    []string{"Hello world"},
			wantEnvNames: []string{
				EnvAxonPrompt,
				EnvAxonModel,
				"ANTHROPIC_API_KEY",
			},
			wantEnvValues: map[string]string{
				EnvAxonPrompt: "Hello world",
				EnvAxonModel:  "",
			},
		},
		{
			name: "Custom image override",
			task: &axonv1alpha1.Task{
				Spec: axonv1alpha1.TaskSpec{
					Type:   AgentTypeClaudeCode,
					Image:  "my-agent:v2",
					Prompt: "Do something",
					Credentials: axonv1alpha1.Credentials{
						Type:      axonv1alpha1.CredentialTypeAPIKey,
						SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
					},
				},
			},
			wantImage:   "my-agent:v2",
			wantCommand: []string{AgentEntrypoint},
			wantArgs:    []string{"Do something"},
			wantEnvNames: []string{
				EnvAxonPrompt,
				EnvAxonModel,
				"ANTHROPIC_API_KEY",
			},
			wantEnvValues: map[string]string{
				EnvAxonPrompt: "Do something",
				EnvAxonModel:  "",
			},
		},
		{
			name: "With model specified",
			task: &axonv1alpha1.Task{
				Spec: axonv1alpha1.TaskSpec{
					Type:   AgentTypeClaudeCode,
					Prompt: "Fix bug",
					Model:  "claude-opus",
					Credentials: axonv1alpha1.Credentials{
						Type:      axonv1alpha1.CredentialTypeOAuth,
						SecretRef: axonv1alpha1.SecretReference{Name: "oauth-secret"},
					},
				},
			},
			wantImage:   ClaudeCodeImage,
			wantCommand: []string{AgentEntrypoint},
			wantArgs:    []string{"Fix bug"},
			wantEnvNames: []string{
				EnvAxonPrompt,
				EnvAxonModel,
				"CLAUDE_CODE_OAUTH_TOKEN",
			},
			wantEnvValues: map[string]string{
				EnvAxonPrompt: "Fix bug",
				EnvAxonModel:  "claude-opus",
			},
		},
		{
			name: "With workspace and secret",
			task: &axonv1alpha1.Task{
				Spec: axonv1alpha1.TaskSpec{
					Type:   AgentTypeClaudeCode,
					Prompt: "Create PR",
					Credentials: axonv1alpha1.Credentials{
						Type:      axonv1alpha1.CredentialTypeAPIKey,
						SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
					},
				},
			},
			workspace: &axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/example/repo.git",
				SecretRef: &axonv1alpha1.SecretReference{
					Name: "gh-token",
				},
			},
			wantImage:   ClaudeCodeImage,
			wantCommand: []string{AgentEntrypoint},
			wantArgs:    []string{"Create PR"},
			wantEnvNames: []string{
				EnvAxonPrompt,
				EnvAxonModel,
				"ANTHROPIC_API_KEY",
				"GITHUB_TOKEN",
				"GH_TOKEN",
			},
			wantEnvValues: map[string]string{
				EnvAxonPrompt: "Create PR",
				EnvAxonModel:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewJobBuilder()
			job, err := builder.Build(tt.task, tt.workspace)
			if err != nil {
				t.Fatalf("Build() error = %v", err)
			}

			container := job.Spec.Template.Spec.Containers[0]

			if container.Image != tt.wantImage {
				t.Errorf("Image = %q, want %q", container.Image, tt.wantImage)
			}

			if len(container.Command) != len(tt.wantCommand) {
				t.Fatalf("Command length = %d, want %d", len(container.Command), len(tt.wantCommand))
			}
			for i, cmd := range container.Command {
				if cmd != tt.wantCommand[i] {
					t.Errorf("Command[%d] = %q, want %q", i, cmd, tt.wantCommand[i])
				}
			}

			if len(container.Args) != len(tt.wantArgs) {
				t.Fatalf("Args length = %d, want %d", len(container.Args), len(tt.wantArgs))
			}
			for i, arg := range container.Args {
				if arg != tt.wantArgs[i] {
					t.Errorf("Args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
				}
			}

			if len(container.Env) != len(tt.wantEnvNames) {
				t.Fatalf("Env length = %d, want %d", len(container.Env), len(tt.wantEnvNames))
			}
			for i, envName := range tt.wantEnvNames {
				if container.Env[i].Name != envName {
					t.Errorf("Env[%d].Name = %q, want %q", i, container.Env[i].Name, envName)
				}
				if expectedValue, ok := tt.wantEnvValues[envName]; ok {
					if container.Env[i].Value != expectedValue {
						t.Errorf("Env[%d].Value = %q, want %q", i, container.Env[i].Value, expectedValue)
					}
				}
			}
		})
	}
}

func TestBuildUnsupportedAgentType(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		Spec: axonv1alpha1.TaskSpec{
			Type:   "unsupported",
			Prompt: "test",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "s"},
			},
		},
	}
	_, err := builder.Build(task, nil)
	if err == nil {
		t.Fatal("Build() expected error for unsupported agent type, got nil")
	}
}
