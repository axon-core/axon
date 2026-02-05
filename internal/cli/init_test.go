package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"init", "--config", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading created file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty config file")
	}
}

func TestInitCommand_CreatesParentDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.yaml")

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"init", "--config", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestInitCommand_RefusesOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"init", "--config", path})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when file exists without --force")
	}
}

func TestInitCommand_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"init", "--config", path, "--force"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) == "existing" {
		t.Fatal("expected file to be overwritten")
	}
}
