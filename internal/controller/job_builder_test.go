package controller

import (
	"strings"
	"testing"

	axonv1alpha1 "github.com/axon-core/axon/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildClaudeCodeJob_DefaultImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-task",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Hello world",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Model: "claude-sonnet-4-20250514",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Default image should be used.
	if container.Image != ClaudeCodeImage {
		t.Errorf("Expected image %q, got %q", ClaudeCodeImage, container.Image)
	}

	// Command should be /axon_entrypoint.sh (uniform interface).
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}

	// Args should be just the prompt.
	if len(container.Args) != 1 || container.Args[0] != "Hello world" {
		t.Errorf("Expected args [Hello world], got %v", container.Args)
	}

	// AXON_MODEL should be set with the correct value.
	foundAxonModel := false
	for _, env := range container.Env {
		if env.Name == "AXON_MODEL" {
			foundAxonModel = true
			if env.Value != "claude-sonnet-4-20250514" {
				t.Errorf("AXON_MODEL value: expected %q, got %q", "claude-sonnet-4-20250514", env.Value)
			}
		}
	}
	if !foundAxonModel {
		t.Error("Expected AXON_MODEL env var to be set")
	}
}

func TestBuildClaudeCodeJob_CustomImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-custom",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Model: "my-model",
			Image: "my-custom-agent:latest",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Custom image should be used.
	if container.Image != "my-custom-agent:latest" {
		t.Errorf("Expected image %q, got %q", "my-custom-agent:latest", container.Image)
	}

	// Command should be /axon_entrypoint.sh (same interface as default).
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}

	// Args should be just the prompt.
	if len(container.Args) != 1 || container.Args[0] != "Fix the bug" {
		t.Errorf("Expected args [Fix the bug], got %v", container.Args)
	}

	// AXON_MODEL should be set with the correct value.
	foundAxonModel := false
	for _, env := range container.Env {
		if env.Name == "AXON_MODEL" {
			foundAxonModel = true
			if env.Value != "my-model" {
				t.Errorf("AXON_MODEL value: expected %q, got %q", "my-model", env.Value)
			}
		}
	}
	if !foundAxonModel {
		t.Error("Expected AXON_MODEL env var to be set")
	}
}

func TestBuildClaudeCodeJob_NoModel(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-no-model",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Hello",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// AXON_MODEL should NOT be set when model is empty.
	for _, env := range container.Env {
		if env.Name == "AXON_MODEL" {
			t.Error("AXON_MODEL should not be set when model is empty")
		}
	}
}

func TestBuildClaudeCodeJob_WorkspaceWithRef(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-workspace",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/example/repo.git",
		Ref:  "main",
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	// Verify git clone args.
	initContainer := job.Spec.Template.Spec.InitContainers[0]
	expectedArgs := []string{
		"clone",
		"--branch", "main", "--no-single-branch", "--depth", "1",
		"--", "https://github.com/example/repo.git", WorkspaceMountPath + "/repo",
	}

	if len(initContainer.Args) != len(expectedArgs) {
		t.Fatalf("Expected %d clone args, got %d: %v", len(expectedArgs), len(initContainer.Args), initContainer.Args)
	}
	for i, arg := range expectedArgs {
		if initContainer.Args[i] != arg {
			t.Errorf("Clone args[%d]: expected %q, got %q", i, arg, initContainer.Args[i])
		}
	}

	// Verify init container runs as ClaudeCodeUID.
	if initContainer.SecurityContext == nil || initContainer.SecurityContext.RunAsUser == nil {
		t.Fatal("Expected init container SecurityContext.RunAsUser to be set")
	}
	if *initContainer.SecurityContext.RunAsUser != ClaudeCodeUID {
		t.Errorf("Expected RunAsUser %d, got %d", ClaudeCodeUID, *initContainer.SecurityContext.RunAsUser)
	}

	// Verify FSGroup.
	if job.Spec.Template.Spec.SecurityContext == nil || job.Spec.Template.Spec.SecurityContext.FSGroup == nil {
		t.Fatal("Expected pod SecurityContext.FSGroup to be set")
	}
	if *job.Spec.Template.Spec.SecurityContext.FSGroup != ClaudeCodeUID {
		t.Errorf("Expected FSGroup %d, got %d", ClaudeCodeUID, *job.Spec.Template.Spec.SecurityContext.FSGroup)
	}

	// Verify main container working dir.
	container := job.Spec.Template.Spec.Containers[0]
	if container.WorkingDir != WorkspaceMountPath+"/repo" {
		t.Errorf("Expected workingDir %q, got %q", WorkspaceMountPath+"/repo", container.WorkingDir)
	}
}

func TestBuildClaudeCodeJob_CustomImageWithWorkspace(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-custom-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Image: "my-agent:v1",
			Model: "gpt-4",
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/example/repo.git",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Custom image with workspace should still use /axon_entrypoint.sh.
	if container.Image != "my-agent:v1" {
		t.Errorf("Expected image %q, got %q", "my-agent:v1", container.Image)
	}
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}
	if len(container.Args) != 1 || container.Args[0] != "Fix the bug" {
		t.Errorf("Expected args [Fix the bug], got %v", container.Args)
	}

	// Should have workspace volume mount and working dir.
	if container.WorkingDir != WorkspaceMountPath+"/repo" {
		t.Errorf("Expected workingDir %q, got %q", WorkspaceMountPath+"/repo", container.WorkingDir)
	}
	if len(container.VolumeMounts) != 1 {
		t.Fatalf("Expected 1 volume mount, got %d", len(container.VolumeMounts))
	}

	// Verify FSGroup.
	if job.Spec.Template.Spec.SecurityContext == nil || job.Spec.Template.Spec.SecurityContext.FSGroup == nil {
		t.Fatal("Expected pod SecurityContext.FSGroup to be set")
	}
	if *job.Spec.Template.Spec.SecurityContext.FSGroup != ClaudeCodeUID {
		t.Errorf("Expected FSGroup %d, got %d", ClaudeCodeUID, *job.Spec.Template.Spec.SecurityContext.FSGroup)
	}

	// Should have AXON_MODEL with correct value, ANTHROPIC_API_KEY, GITHUB_TOKEN, GH_TOKEN.
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		} else {
			envMap[env.Name] = "(from-secret)"
		}
	}
	for _, name := range []string{"AXON_MODEL", "ANTHROPIC_API_KEY", "GITHUB_TOKEN", "GH_TOKEN"} {
		if _, ok := envMap[name]; !ok {
			t.Errorf("Expected env var %q to be set", name)
		}
	}
	if envMap["AXON_MODEL"] != "gpt-4" {
		t.Errorf("AXON_MODEL value: expected %q, got %q", "gpt-4", envMap["AXON_MODEL"])
	}
}

func TestBuildClaudeCodeJob_WorkspaceWithSecretRefPersistsCredentialHelper(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-persist-cred",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/example/repo.git",
		Ref:  "main",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	initContainer := job.Spec.Template.Spec.InitContainers[0]

	// Verify the init container command uses sh -c.
	if len(initContainer.Command) != 3 || initContainer.Command[0] != "sh" || initContainer.Command[1] != "-c" {
		t.Fatalf("Expected command [sh -c ...], got %v", initContainer.Command)
	}

	script := initContainer.Command[2]

	// The script must clone with an inline credential helper AND persist it
	// to the repo config so the agent container can authenticate with git.
	if !strings.Contains(script, "git -c credential.helper=") {
		t.Error("Expected init container script to include inline credential helper for clone")
	}
	if !strings.Contains(script, "git -C "+WorkspaceMountPath+"/repo config credential.helper") {
		t.Error("Expected init container script to persist credential helper in repo config")
	}
}

func TestBuildClaudeCodeJob_EnterpriseWorkspaceSetsGHHostAndEnterpriseToken(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ghe",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.example.com/my-org/my-repo.git",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		} else {
			envMap[env.Name] = "(from-secret)"
		}
	}

	// GH_HOST should be set for enterprise.
	if envMap["GH_HOST"] != "github.example.com" {
		t.Errorf("Expected GH_HOST = %q, got %q", "github.example.com", envMap["GH_HOST"])
	}
	// GH_ENTERPRISE_TOKEN should be set instead of GH_TOKEN for enterprise hosts.
	if _, ok := envMap["GH_ENTERPRISE_TOKEN"]; !ok {
		t.Error("Expected GH_ENTERPRISE_TOKEN to be set for enterprise workspace")
	}
	if _, ok := envMap["GH_TOKEN"]; ok {
		t.Error("GH_TOKEN should not be set for enterprise workspace")
	}
	// GITHUB_TOKEN should still be set (used for git credential helper).
	if _, ok := envMap["GITHUB_TOKEN"]; !ok {
		t.Error("Expected GITHUB_TOKEN to be set for enterprise workspace")
	}

	initContainer := job.Spec.Template.Spec.InitContainers[0]
	initEnvMap := map[string]string{}
	for _, env := range initContainer.Env {
		if env.Value != "" {
			initEnvMap[env.Name] = env.Value
		} else {
			initEnvMap[env.Name] = "(from-secret)"
		}
	}
	if initEnvMap["GH_HOST"] != "github.example.com" {
		t.Errorf("Expected init container GH_HOST = %q, got %q", "github.example.com", initEnvMap["GH_HOST"])
	}
	if _, ok := initEnvMap["GH_ENTERPRISE_TOKEN"]; !ok {
		t.Error("Expected GH_ENTERPRISE_TOKEN in init container for enterprise workspace")
	}
	if _, ok := initEnvMap["GH_TOKEN"]; ok {
		t.Error("GH_TOKEN should not be set in init container for enterprise workspace")
	}
}

func TestBuildClaudeCodeJob_GithubComWorkspaceUsesGHToken(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-no-ghe",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/my-org/my-repo.git",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		} else {
			envMap[env.Name] = "(from-secret)"
		}
	}

	// GH_HOST should NOT be set for github.com.
	if _, ok := envMap["GH_HOST"]; ok {
		t.Error("GH_HOST should not be set for github.com workspace")
	}
	// GH_TOKEN should be set for github.com.
	if _, ok := envMap["GH_TOKEN"]; !ok {
		t.Error("Expected GH_TOKEN to be set for github.com workspace")
	}
	// GH_ENTERPRISE_TOKEN should NOT be set for github.com.
	if _, ok := envMap["GH_ENTERPRISE_TOKEN"]; ok {
		t.Error("GH_ENTERPRISE_TOKEN should not be set for github.com workspace")
	}
}

func TestBuildCodexJob_DefaultImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-codex",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeCodex,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "openai-secret"},
			},
			Model: "gpt-4.1",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Default codex image should be used.
	if container.Image != CodexImage {
		t.Errorf("Expected image %q, got %q", CodexImage, container.Image)
	}

	// Container name should match the agent type.
	if container.Name != AgentTypeCodex {
		t.Errorf("Expected container name %q, got %q", AgentTypeCodex, container.Name)
	}

	// Command should be /axon_entrypoint.sh (uniform interface).
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}

	// Args should be just the prompt.
	if len(container.Args) != 1 || container.Args[0] != "Fix the bug" {
		t.Errorf("Expected args [Fix the bug], got %v", container.Args)
	}

	// AXON_MODEL should be set.
	foundAxonModel := false
	for _, env := range container.Env {
		if env.Name == "AXON_MODEL" {
			foundAxonModel = true
			if env.Value != "gpt-4.1" {
				t.Errorf("AXON_MODEL value: expected %q, got %q", "gpt-4.1", env.Value)
			}
		}
	}
	if !foundAxonModel {
		t.Error("Expected AXON_MODEL env var to be set")
	}

	// CODEX_API_KEY should be set (not ANTHROPIC_API_KEY).
	foundCodexKey := false
	for _, env := range container.Env {
		if env.Name == "CODEX_API_KEY" {
			foundCodexKey = true
			if env.ValueFrom == nil || env.ValueFrom.SecretKeyRef == nil {
				t.Error("Expected CODEX_API_KEY to reference a secret")
			} else {
				if env.ValueFrom.SecretKeyRef.Name != "openai-secret" {
					t.Errorf("Expected secret name %q, got %q", "openai-secret", env.ValueFrom.SecretKeyRef.Name)
				}
				if env.ValueFrom.SecretKeyRef.Key != "CODEX_API_KEY" {
					t.Errorf("Expected secret key %q, got %q", "CODEX_API_KEY", env.ValueFrom.SecretKeyRef.Key)
				}
			}
		}
		if env.Name == "ANTHROPIC_API_KEY" {
			t.Error("ANTHROPIC_API_KEY should not be set for codex agent type")
		}
	}
	if !foundCodexKey {
		t.Error("Expected CODEX_API_KEY env var to be set")
	}
}

func TestBuildCodexJob_CustomImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-codex-custom",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeCodex,
			Prompt: "Refactor the module",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "openai-secret"},
			},
			Image: "my-codex:v2",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Custom image should be used.
	if container.Image != "my-codex:v2" {
		t.Errorf("Expected image %q, got %q", "my-codex:v2", container.Image)
	}

	// Command should be /axon_entrypoint.sh.
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}
}

func TestBuildCodexJob_WithWorkspace(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-codex-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeCodex,
			Prompt: "Fix the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "openai-secret"},
			},
			Model: "gpt-4.1",
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/example/repo.git",
		Ref:  "main",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Should have workspace volume mount and working dir.
	if container.WorkingDir != WorkspaceMountPath+"/repo" {
		t.Errorf("Expected workingDir %q, got %q", WorkspaceMountPath+"/repo", container.WorkingDir)
	}
	if len(container.VolumeMounts) != 1 {
		t.Fatalf("Expected 1 volume mount, got %d", len(container.VolumeMounts))
	}

	// Should have CODEX_API_KEY (not ANTHROPIC_API_KEY), AXON_MODEL, GITHUB_TOKEN, GH_TOKEN.
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		} else {
			envMap[env.Name] = "(from-secret)"
		}
	}
	for _, name := range []string{"AXON_MODEL", "CODEX_API_KEY", "GITHUB_TOKEN", "GH_TOKEN"} {
		if _, ok := envMap[name]; !ok {
			t.Errorf("Expected env var %q to be set", name)
		}
	}
	if _, ok := envMap["ANTHROPIC_API_KEY"]; ok {
		t.Error("ANTHROPIC_API_KEY should not be set for codex agent type")
	}

	// Verify init container and FSGroup.
	if len(job.Spec.Template.Spec.InitContainers) != 1 {
		t.Fatalf("Expected 1 init container, got %d", len(job.Spec.Template.Spec.InitContainers))
	}
	initContainer := job.Spec.Template.Spec.InitContainers[0]
	if initContainer.SecurityContext == nil || initContainer.SecurityContext.RunAsUser == nil {
		t.Fatal("Expected init container SecurityContext.RunAsUser to be set")
	}
	if *initContainer.SecurityContext.RunAsUser != AgentUID {
		t.Errorf("Expected RunAsUser %d, got %d", AgentUID, *initContainer.SecurityContext.RunAsUser)
	}
}

func TestBuildCodexJob_OAuthCredentials(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-codex-oauth",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeCodex,
			Prompt: "Review the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeOAuth,
				SecretRef: axonv1alpha1.SecretReference{Name: "codex-oauth"},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// CODEX_API_KEY should be set for codex oauth.
	foundCodexKey := false
	for _, env := range container.Env {
		if env.Name == "CODEX_API_KEY" {
			foundCodexKey = true
			if env.ValueFrom == nil || env.ValueFrom.SecretKeyRef == nil {
				t.Error("Expected CODEX_API_KEY to reference a secret")
			} else {
				if env.ValueFrom.SecretKeyRef.Name != "codex-oauth" {
					t.Errorf("Expected secret name %q, got %q", "codex-oauth", env.ValueFrom.SecretKeyRef.Name)
				}
				if env.ValueFrom.SecretKeyRef.Key != "CODEX_API_KEY" {
					t.Errorf("Expected secret key %q, got %q", "CODEX_API_KEY", env.ValueFrom.SecretKeyRef.Key)
				}
			}
		}
		if env.Name == "CLAUDE_CODE_OAUTH_TOKEN" {
			t.Error("CLAUDE_CODE_OAUTH_TOKEN should not be set for codex agent type")
		}
	}
	if !foundCodexKey {
		t.Error("Expected CODEX_API_KEY env var to be set")
	}
}

func TestBuildGeminiJob_DefaultImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gemini",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeGemini,
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "gemini-secret"},
			},
			Model: "gemini-2.5-pro",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Default gemini image should be used.
	if container.Image != GeminiImage {
		t.Errorf("Expected image %q, got %q", GeminiImage, container.Image)
	}

	// Container name should match the agent type.
	if container.Name != AgentTypeGemini {
		t.Errorf("Expected container name %q, got %q", AgentTypeGemini, container.Name)
	}

	// Command should be /axon_entrypoint.sh (uniform interface).
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}

	// Args should be just the prompt.
	if len(container.Args) != 1 || container.Args[0] != "Fix the bug" {
		t.Errorf("Expected args [Fix the bug], got %v", container.Args)
	}

	// AXON_MODEL should be set.
	foundAxonModel := false
	for _, env := range container.Env {
		if env.Name == "AXON_MODEL" {
			foundAxonModel = true
			if env.Value != "gemini-2.5-pro" {
				t.Errorf("AXON_MODEL value: expected %q, got %q", "gemini-2.5-pro", env.Value)
			}
		}
	}
	if !foundAxonModel {
		t.Error("Expected AXON_MODEL env var to be set")
	}

	// GEMINI_API_KEY should be set (not ANTHROPIC_API_KEY or CODEX_API_KEY).
	foundGeminiKey := false
	for _, env := range container.Env {
		if env.Name == "GEMINI_API_KEY" {
			foundGeminiKey = true
			if env.ValueFrom == nil || env.ValueFrom.SecretKeyRef == nil {
				t.Error("Expected GEMINI_API_KEY to reference a secret")
			} else {
				if env.ValueFrom.SecretKeyRef.Name != "gemini-secret" {
					t.Errorf("Expected secret name %q, got %q", "gemini-secret", env.ValueFrom.SecretKeyRef.Name)
				}
				if env.ValueFrom.SecretKeyRef.Key != "GEMINI_API_KEY" {
					t.Errorf("Expected secret key %q, got %q", "GEMINI_API_KEY", env.ValueFrom.SecretKeyRef.Key)
				}
			}
		}
		if env.Name == "ANTHROPIC_API_KEY" {
			t.Error("ANTHROPIC_API_KEY should not be set for gemini agent type")
		}
		if env.Name == "CODEX_API_KEY" {
			t.Error("CODEX_API_KEY should not be set for gemini agent type")
		}
	}
	if !foundGeminiKey {
		t.Error("Expected GEMINI_API_KEY env var to be set")
	}
}

func TestBuildGeminiJob_CustomImage(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gemini-custom",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeGemini,
			Prompt: "Refactor the module",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "gemini-secret"},
			},
			Image: "my-gemini:v2",
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Custom image should be used.
	if container.Image != "my-gemini:v2" {
		t.Errorf("Expected image %q, got %q", "my-gemini:v2", container.Image)
	}

	// Command should be /axon_entrypoint.sh.
	if len(container.Command) != 1 || container.Command[0] != "/axon_entrypoint.sh" {
		t.Errorf("Expected command [/axon_entrypoint.sh], got %v", container.Command)
	}
}

func TestBuildGeminiJob_WithWorkspace(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gemini-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeGemini,
			Prompt: "Fix the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "gemini-secret"},
			},
			Model: "gemini-2.5-pro",
		},
	}

	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/example/repo.git",
		Ref:  "main",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	job, err := builder.Build(task, workspace)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Should have workspace volume mount and working dir.
	if container.WorkingDir != WorkspaceMountPath+"/repo" {
		t.Errorf("Expected workingDir %q, got %q", WorkspaceMountPath+"/repo", container.WorkingDir)
	}
	if len(container.VolumeMounts) != 1 {
		t.Fatalf("Expected 1 volume mount, got %d", len(container.VolumeMounts))
	}

	// Should have GEMINI_API_KEY (not ANTHROPIC_API_KEY), AXON_MODEL, GITHUB_TOKEN, GH_TOKEN.
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		} else {
			envMap[env.Name] = "(from-secret)"
		}
	}
	for _, name := range []string{"AXON_MODEL", "GEMINI_API_KEY", "GITHUB_TOKEN", "GH_TOKEN"} {
		if _, ok := envMap[name]; !ok {
			t.Errorf("Expected env var %q to be set", name)
		}
	}
	if _, ok := envMap["ANTHROPIC_API_KEY"]; ok {
		t.Error("ANTHROPIC_API_KEY should not be set for gemini agent type")
	}
	if _, ok := envMap["CODEX_API_KEY"]; ok {
		t.Error("CODEX_API_KEY should not be set for gemini agent type")
	}

	// Verify init container and FSGroup.
	if len(job.Spec.Template.Spec.InitContainers) != 1 {
		t.Fatalf("Expected 1 init container, got %d", len(job.Spec.Template.Spec.InitContainers))
	}
	initContainer := job.Spec.Template.Spec.InitContainers[0]
	if initContainer.SecurityContext == nil || initContainer.SecurityContext.RunAsUser == nil {
		t.Fatal("Expected init container SecurityContext.RunAsUser to be set")
	}
	if *initContainer.SecurityContext.RunAsUser != AgentUID {
		t.Errorf("Expected RunAsUser %d, got %d", AgentUID, *initContainer.SecurityContext.RunAsUser)
	}
}

func TestBuildGeminiJob_OAuthCredentials(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-gemini-oauth",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeGemini,
			Prompt: "Review the code",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeOAuth,
				SecretRef: axonv1alpha1.SecretReference{Name: "gemini-oauth"},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// GEMINI_API_KEY should be set for gemini oauth.
	foundGeminiKey := false
	for _, env := range container.Env {
		if env.Name == "GEMINI_API_KEY" {
			foundGeminiKey = true
			if env.ValueFrom == nil || env.ValueFrom.SecretKeyRef == nil {
				t.Error("Expected GEMINI_API_KEY to reference a secret")
			} else {
				if env.ValueFrom.SecretKeyRef.Name != "gemini-oauth" {
					t.Errorf("Expected secret name %q, got %q", "gemini-oauth", env.ValueFrom.SecretKeyRef.Name)
				}
				if env.ValueFrom.SecretKeyRef.Key != "GEMINI_API_KEY" {
					t.Errorf("Expected secret key %q, got %q", "GEMINI_API_KEY", env.ValueFrom.SecretKeyRef.Key)
				}
			}
		}
		if env.Name == "CLAUDE_CODE_OAUTH_TOKEN" {
			t.Error("CLAUDE_CODE_OAUTH_TOKEN should not be set for gemini agent type")
		}
		if env.Name == "CODEX_API_KEY" {
			t.Error("CODEX_API_KEY should not be set for gemini agent type")
		}
	}
	if !foundGeminiKey {
		t.Error("Expected GEMINI_API_KEY env var to be set")
	}
}

func TestBuildClaudeCodeJob_UnsupportedType(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-unsupported",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   "unsupported-agent",
			Prompt: "Hello",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	_, err := builder.Build(task, nil)
	if err == nil {
		t.Fatal("Expected error for unsupported agent type, got nil")
	}
}

func int64Ptr(v int64) *int64 { return &v }

func TestBuildJob_PodOverridesResources(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-resources",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			PodOverrides: &axonv1alpha1.PodOverrides{
				Resources: &corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse("512Mi"),
						corev1.ResourceCPU:    resource.MustParse("500m"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse("2Gi"),
						corev1.ResourceCPU:    resource.MustParse("2"),
					},
				},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	memReq := container.Resources.Requests[corev1.ResourceMemory]
	if memReq.String() != "512Mi" {
		t.Errorf("Expected memory request 512Mi, got %s", memReq.String())
	}
	cpuReq := container.Resources.Requests[corev1.ResourceCPU]
	if cpuReq.String() != "500m" {
		t.Errorf("Expected CPU request 500m, got %s", cpuReq.String())
	}
	memLimit := container.Resources.Limits[corev1.ResourceMemory]
	if memLimit.String() != "2Gi" {
		t.Errorf("Expected memory limit 2Gi, got %s", memLimit.String())
	}
	cpuLimit := container.Resources.Limits[corev1.ResourceCPU]
	if cpuLimit.String() != "2" {
		t.Errorf("Expected CPU limit 2, got %s", cpuLimit.String())
	}
}

func TestBuildJob_PodOverridesActiveDeadlineSeconds(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deadline",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			PodOverrides: &axonv1alpha1.PodOverrides{
				ActiveDeadlineSeconds: int64Ptr(1800),
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	if job.Spec.ActiveDeadlineSeconds == nil {
		t.Fatal("Expected ActiveDeadlineSeconds to be set")
	}
	if *job.Spec.ActiveDeadlineSeconds != 1800 {
		t.Errorf("Expected ActiveDeadlineSeconds 1800, got %d", *job.Spec.ActiveDeadlineSeconds)
	}
}

func TestBuildJob_PodOverridesEnv(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-env",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Model: "claude-sonnet-4-20250514",
			PodOverrides: &axonv1alpha1.PodOverrides{
				Env: []corev1.EnvVar{
					{Name: "HTTP_PROXY", Value: "http://proxy:8080"},
					{Name: "NO_PROXY", Value: "localhost"},
				},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		}
	}

	// User env vars should be present.
	if envMap["HTTP_PROXY"] != "http://proxy:8080" {
		t.Errorf("Expected HTTP_PROXY=http://proxy:8080, got %q", envMap["HTTP_PROXY"])
	}
	if envMap["NO_PROXY"] != "localhost" {
		t.Errorf("Expected NO_PROXY=localhost, got %q", envMap["NO_PROXY"])
	}

	// Built-in env vars should still be present.
	if envMap["AXON_MODEL"] != "claude-sonnet-4-20250514" {
		t.Errorf("Expected AXON_MODEL to still be set, got %q", envMap["AXON_MODEL"])
	}
}

func TestBuildJob_PodOverridesEnvBuiltinPrecedence(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-env-precedence",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Model: "claude-sonnet-4-20250514",
			PodOverrides: &axonv1alpha1.PodOverrides{
				Env: []corev1.EnvVar{
					// Attempt to override a built-in env var.
					{Name: "AXON_MODEL", Value: "should-not-take-effect"},
				},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// User env vars that collide with built-in names should be filtered out
	// so that built-in vars always take precedence.
	var axonModelCount int
	for _, e := range container.Env {
		if e.Name == "AXON_MODEL" {
			axonModelCount++
			if e.Value != "claude-sonnet-4-20250514" {
				t.Errorf("Expected AXON_MODEL value %q, got %q", "claude-sonnet-4-20250514", e.Value)
			}
		}
	}
	if axonModelCount != 1 {
		t.Errorf("Expected exactly 1 AXON_MODEL env var, got %d", axonModelCount)
	}
}

func TestBuildJob_PodOverridesNodeSelector(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-selector",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			PodOverrides: &axonv1alpha1.PodOverrides{
				NodeSelector: map[string]string{
					"workload-type": "ai-agent",
					"gpu":           "true",
				},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	ns := job.Spec.Template.Spec.NodeSelector
	if ns == nil {
		t.Fatal("Expected NodeSelector to be set")
	}
	if ns["workload-type"] != "ai-agent" {
		t.Errorf("Expected nodeSelector workload-type=ai-agent, got %q", ns["workload-type"])
	}
	if ns["gpu"] != "true" {
		t.Errorf("Expected nodeSelector gpu=true, got %q", ns["gpu"])
	}
}

func TestBuildJob_PodOverridesAllFields(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-all-overrides",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeCodex,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "openai-secret"},
			},
			PodOverrides: &axonv1alpha1.PodOverrides{
				Resources: &corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse("4Gi"),
					},
				},
				ActiveDeadlineSeconds: int64Ptr(3600),
				Env: []corev1.EnvVar{
					{Name: "HTTPS_PROXY", Value: "http://proxy:8080"},
				},
				NodeSelector: map[string]string{
					"pool": "agents",
				},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// Resources
	memLimit := container.Resources.Limits[corev1.ResourceMemory]
	if memLimit.String() != "4Gi" {
		t.Errorf("Expected memory limit 4Gi, got %s", memLimit.String())
	}

	// ActiveDeadlineSeconds
	if job.Spec.ActiveDeadlineSeconds == nil || *job.Spec.ActiveDeadlineSeconds != 3600 {
		t.Errorf("Expected ActiveDeadlineSeconds 3600, got %v", job.Spec.ActiveDeadlineSeconds)
	}

	// Env
	envMap := map[string]string{}
	for _, env := range container.Env {
		if env.Value != "" {
			envMap[env.Name] = env.Value
		}
	}
	if envMap["HTTPS_PROXY"] != "http://proxy:8080" {
		t.Errorf("Expected HTTPS_PROXY=http://proxy:8080, got %q", envMap["HTTPS_PROXY"])
	}

	// NodeSelector
	if job.Spec.Template.Spec.NodeSelector["pool"] != "agents" {
		t.Errorf("Expected nodeSelector pool=agents, got %q", job.Spec.Template.Spec.NodeSelector["pool"])
	}
}

func TestBuildJob_NoPodOverrides(t *testing.T) {
	builder := NewJobBuilder()
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-no-overrides",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   AgentTypeClaudeCode,
			Prompt: "Fix issue",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
		},
	}

	job, err := builder.Build(task, nil)
	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	container := job.Spec.Template.Spec.Containers[0]

	// No resources should be set.
	if len(container.Resources.Requests) != 0 || len(container.Resources.Limits) != 0 {
		t.Error("Expected no resources to be set when PodOverrides is nil")
	}

	// No ActiveDeadlineSeconds.
	if job.Spec.ActiveDeadlineSeconds != nil {
		t.Error("Expected no ActiveDeadlineSeconds when PodOverrides is nil")
	}

	// No NodeSelector.
	if job.Spec.Template.Spec.NodeSelector != nil {
		t.Error("Expected no NodeSelector when PodOverrides is nil")
	}
}
