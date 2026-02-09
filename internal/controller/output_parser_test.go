package controller

import (
	"testing"
)

func TestParseOutputs(t *testing.T) {
	tests := []struct {
		name     string
		logData  string
		expected []string
	}{
		{
			name:     "no markers",
			logData:  "some random log output\nmore logs\n",
			expected: nil,
		},
		{
			name:     "empty between markers",
			logData:  "---AXON_OUTPUTS_START---\n---AXON_OUTPUTS_END---\n",
			expected: nil,
		},
		{
			name:     "whitespace only between markers",
			logData:  "---AXON_OUTPUTS_START---\n  \n  \n---AXON_OUTPUTS_END---\n",
			expected: nil,
		},
		{
			name:     "branch only",
			logData:  "---AXON_OUTPUTS_START---\nbranch: axon-task-123\n---AXON_OUTPUTS_END---\n",
			expected: []string{"branch: axon-task-123"},
		},
		{
			name: "branch and PR URL",
			logData: "---AXON_OUTPUTS_START---\nbranch: axon-task-123\n" +
				"https://github.com/org/repo/pull/456\n---AXON_OUTPUTS_END---\n",
			expected: []string{
				"branch: axon-task-123",
				"https://github.com/org/repo/pull/456",
			},
		},
		{
			name: "branch and multiple PR URLs",
			logData: "---AXON_OUTPUTS_START---\nbranch: feature-branch\n" +
				"https://github.com/org/repo/pull/1\n" +
				"https://github.com/org/repo/pull/2\n---AXON_OUTPUTS_END---\n",
			expected: []string{
				"branch: feature-branch",
				"https://github.com/org/repo/pull/1",
				"https://github.com/org/repo/pull/2",
			},
		},
		{
			name: "markers in noisy log data",
			logData: "Starting agent...\nProcessing task...\nDone.\n" +
				"---AXON_OUTPUTS_START---\nbranch: my-branch\n" +
				"https://github.com/org/repo/pull/99\n---AXON_OUTPUTS_END---\n" +
				"Exiting with code 0\n",
			expected: []string{
				"branch: my-branch",
				"https://github.com/org/repo/pull/99",
			},
		},
		{
			name:     "start marker without end marker",
			logData:  "---AXON_OUTPUTS_START---\nbranch: broken\n",
			expected: nil,
		},
		{
			name:     "end marker before start marker",
			logData:  "---AXON_OUTPUTS_END---\n---AXON_OUTPUTS_START---\nbranch: wrong-order\n",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseOutputs(tt.logData)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d items, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("item %d: expected %q, got %q", i, v, result[i])
				}
			}
		})
	}
}
