package source

import (
	"testing"
)

func TestSortByLabelPriority_BasicOrdering(t *testing.T) {
	items := []WorkItem{
		{ID: "1", Labels: []string{"priority/backlog"}},
		{ID: "2", Labels: []string{"priority/critical-urgent"}},
		{ID: "3", Labels: []string{"priority/imporant-soon"}},
	}

	priorityLabels := []string{
		"priority/critical-urgent",
		"priority/imporant-soon",
		"priority/import-longterm",
		"priority/backlog",
	}

	SortByLabelPriority(items, priorityLabels)

	expected := []string{"2", "3", "1"}
	for i, want := range expected {
		if items[i].ID != want {
			t.Errorf("Position %d: got ID %q, want %q", i, items[i].ID, want)
		}
	}
}

func TestSortByLabelPriority_NoMatchingLabelsGoLast(t *testing.T) {
	items := []WorkItem{
		{ID: "1", Labels: []string{"unrelated-label"}},
		{ID: "2", Labels: []string{"priority/critical-urgent"}},
		{ID: "3", Labels: nil},
	}

	priorityLabels := []string{"priority/critical-urgent", "priority/backlog"}

	SortByLabelPriority(items, priorityLabels)

	if items[0].ID != "2" {
		t.Errorf("Expected critical-urgent item first, got ID %q", items[0].ID)
	}
	// Items without matching labels should be last, maintaining relative order
	if items[1].ID != "1" || items[2].ID != "3" {
		t.Errorf("Expected unmatched items last in original order, got %q, %q", items[1].ID, items[2].ID)
	}
}

func TestSortByLabelPriority_StableSort(t *testing.T) {
	items := []WorkItem{
		{ID: "1", Labels: []string{"priority/backlog"}},
		{ID: "2", Labels: []string{"priority/backlog"}},
		{ID: "3", Labels: []string{"priority/backlog"}},
	}

	priorityLabels := []string{"priority/critical-urgent", "priority/backlog"}

	SortByLabelPriority(items, priorityLabels)

	// Same priority items should retain original order
	expected := []string{"1", "2", "3"}
	for i, want := range expected {
		if items[i].ID != want {
			t.Errorf("Position %d: got ID %q, want %q (stable sort broken)", i, items[i].ID, want)
		}
	}
}

func TestSortByLabelPriority_EmptyPriorityLabels(t *testing.T) {
	items := []WorkItem{
		{ID: "1", Labels: []string{"priority/backlog"}},
		{ID: "2", Labels: []string{"priority/critical-urgent"}},
	}

	SortByLabelPriority(items, nil)

	// Order should not change
	if items[0].ID != "1" || items[1].ID != "2" {
		t.Errorf("Expected original order when priorityLabels is nil, got %q, %q", items[0].ID, items[1].ID)
	}

	SortByLabelPriority(items, []string{})

	// Order should not change
	if items[0].ID != "1" || items[1].ID != "2" {
		t.Errorf("Expected original order when priorityLabels is empty, got %q, %q", items[0].ID, items[1].ID)
	}
}

func TestSortByLabelPriority_EmptyItems(t *testing.T) {
	// Should not panic
	SortByLabelPriority(nil, []string{"priority/critical-urgent"})
	SortByLabelPriority([]WorkItem{}, []string{"priority/critical-urgent"})
}

func TestSortByLabelPriority_MultipleLabelsUseBestMatch(t *testing.T) {
	items := []WorkItem{
		{ID: "1", Labels: []string{"kind/bug", "priority/backlog"}},
		{ID: "2", Labels: []string{"kind/feature", "priority/critical-urgent", "priority/backlog"}},
		{ID: "3", Labels: []string{"priority/imporant-soon", "kind/bug"}},
	}

	priorityLabels := []string{
		"priority/critical-urgent",
		"priority/imporant-soon",
		"priority/backlog",
	}

	SortByLabelPriority(items, priorityLabels)

	expected := []string{"2", "3", "1"}
	for i, want := range expected {
		if items[i].ID != want {
			t.Errorf("Position %d: got ID %q, want %q", i, items[i].ID, want)
		}
	}
}

func TestSortByLabelPriority_IssueScenario(t *testing.T) {
	// Reproduce the exact scenario from the issue
	items := []WorkItem{
		{ID: "500", Labels: []string{"priority/import-longterm"}},
		{ID: "400", Labels: []string{"priority/imporant-soon"}},
		{ID: "300", Labels: []string{"priority/backlog"}},
		{ID: "200", Labels: []string{"priority/backlog"}},
		{ID: "100", Labels: []string{"priority/critical-urgent"}},
	}

	priorityLabels := []string{
		"priority/critical-urgent",
		"priority/imporant-soon",
		"priority/import-longterm",
		"priority/backlog",
	}

	SortByLabelPriority(items, priorityLabels)

	expected := []string{"100", "400", "500", "300", "200"}
	for i, want := range expected {
		if items[i].ID != want {
			t.Errorf("Position %d: got ID %q, want %q", i, items[i].ID, want)
		}
	}
}

func TestLabelPriorityIndex(t *testing.T) {
	priorityLabels := []string{"priority/critical-urgent", "priority/imporant-soon", "priority/backlog"}

	tests := []struct {
		name       string
		itemLabels []string
		want       int
	}{
		{"critical", []string{"priority/critical-urgent"}, 0},
		{"important", []string{"priority/imporant-soon"}, 1},
		{"backlog", []string{"priority/backlog"}, 2},
		{"no match", []string{"unrelated"}, 3},
		{"no labels", nil, 3},
		{"best match wins", []string{"priority/backlog", "priority/critical-urgent"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := labelPriorityIndex(tt.itemLabels, priorityLabels)
			if got != tt.want {
				t.Errorf("labelPriorityIndex(%v) = %d, want %d", tt.itemLabels, got, tt.want)
			}
		})
	}
}
