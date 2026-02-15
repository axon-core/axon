package cli

import (
	"strings"
	"testing"
)

func TestSuspendCommand_MissingName(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"suspend", "taskspawner"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when name is missing")
	}
	if !strings.Contains(err.Error(), "task spawner name is required") {
		t.Errorf("Expected 'task spawner name is required' error, got: %v", err)
	}
}

func TestSuspendCommand_TooManyArgs(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"suspend", "taskspawner", "a", "b"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error with too many arguments")
	}
	if !strings.Contains(err.Error(), "too many arguments") {
		t.Errorf("Expected 'too many arguments' error, got: %v", err)
	}
}

func TestSuspendCommand_NoResourceType(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"suspend"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when no resource type specified")
	}
	if !strings.Contains(err.Error(), "must specify a resource type") {
		t.Errorf("Expected 'must specify a resource type' error, got: %v", err)
	}
}

func TestResumeCommand_MissingName(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"resume", "taskspawner"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when name is missing")
	}
	if !strings.Contains(err.Error(), "task spawner name is required") {
		t.Errorf("Expected 'task spawner name is required' error, got: %v", err)
	}
}

func TestResumeCommand_TooManyArgs(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"resume", "taskspawner", "a", "b"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error with too many arguments")
	}
	if !strings.Contains(err.Error(), "too many arguments") {
		t.Errorf("Expected 'too many arguments' error, got: %v", err)
	}
}

func TestResumeCommand_NoResourceType(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"resume"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error when no resource type specified")
	}
	if !strings.Contains(err.Error(), "must specify a resource type") {
		t.Errorf("Expected 'must specify a resource type' error, got: %v", err)
	}
}
