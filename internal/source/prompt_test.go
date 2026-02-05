package source

import (
	"strings"
	"testing"
)

func TestRenderPromptDefault(t *testing.T) {
	item := WorkItem{
		ID:     "42",
		Number: 42,
		Title:  "Fix the login bug",
		Body:   "Users cannot log in after the update.",
		URL:    "https://github.com/o/r/issues/42",
		Labels: []string{"bug"},
	}

	result, err := RenderPrompt("", item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "#42") {
		t.Errorf("expected issue number in output: %s", result)
	}
	if !strings.Contains(result, "Fix the login bug") {
		t.Errorf("expected title in output: %s", result)
	}
	if !strings.Contains(result, "Users cannot log in") {
		t.Errorf("expected body in output: %s", result)
	}
}

func TestRenderPromptCustom(t *testing.T) {
	item := WorkItem{
		Number: 7,
		Title:  "Add dark mode",
	}

	result, err := RenderPrompt("Fix {{.Title}} (issue #{{.Number}})", item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Fix Add dark mode (issue #7)"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderPromptWithComments(t *testing.T) {
	item := WorkItem{
		Number:   1,
		Title:    "Test",
		Body:     "Body",
		Comments: "A comment",
	}

	result, err := RenderPrompt("", item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Comments:") {
		t.Errorf("expected Comments section in output: %s", result)
	}
	if !strings.Contains(result, "A comment") {
		t.Errorf("expected comment text in output: %s", result)
	}

	// Without comments
	item.Comments = ""
	result, err = RenderPrompt("", item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "Comments:") {
		t.Errorf("expected no Comments section in output: %s", result)
	}
}

func TestRenderPromptAllVariables(t *testing.T) {
	item := WorkItem{
		ID:       "99",
		Number:   99,
		Title:    "T",
		Body:     "B",
		URL:      "U",
		Labels:   []string{"a", "b"},
		Comments: "C",
	}

	tmpl := "{{.ID}} {{.Number}} {{.Title}} {{.Body}} {{.URL}} {{.Labels}} {{.Comments}}"
	result, err := RenderPrompt(tmpl, item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "99 99 T B U a, b C"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderPromptInvalidTemplate(t *testing.T) {
	item := WorkItem{}

	_, err := RenderPrompt("{{.Invalid", item)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}
