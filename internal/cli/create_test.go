package cli

import (
	"testing"
)

func TestCreateCommand_RequiresSubcommand(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no resource type specified")
	}
}

func TestCreateWorkspaceCommand_RequiresName(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "workspace", "--repo", "https://github.com/example/repo.git"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --name not specified")
	}
}

func TestCreateWorkspaceCommand_RequiresRepo(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "workspace", "--name", "my-workspace"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --repo not specified")
	}
}

func TestCreateTaskSpawnerCommand_RequiresName(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "taskspawner", "--workspace", "my-workspace"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --name not specified")
	}
}

func TestCreateTaskSpawnerCommand_RequiresWorkspace(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"create", "taskspawner", "--name", "my-spawner"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --workspace not specified")
	}
}
