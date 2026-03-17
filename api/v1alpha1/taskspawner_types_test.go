package v1alpha1

import "testing"

func TestEffectiveOn(t *testing.T) {
	tests := []struct {
		name string
		spec TaskSpawnerSpec
		want *On
	}{
		{
			name: "returns On when On is set",
			spec: TaskSpawnerSpec{
				On: &On{GitHubIssues: &GitHubIssues{}},
			},
			want: &On{GitHubIssues: &GitHubIssues{}},
		},
		{
			name: "returns When when only When is set",
			spec: TaskSpawnerSpec{
				When: &On{Cron: &Cron{Schedule: "0 9 * * 1"}},
			},
			want: &On{Cron: &Cron{Schedule: "0 9 * * 1"}},
		},
		{
			name: "prefers On over When when both are set",
			spec: TaskSpawnerSpec{
				On:   &On{GitHubIssues: &GitHubIssues{}},
				When: &On{Cron: &Cron{Schedule: "0 9 * * 1"}},
			},
			want: &On{GitHubIssues: &GitHubIssues{}},
		},
		{
			name: "returns nil when neither is set",
			spec: TaskSpawnerSpec{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.EffectiveOn()
			if tt.want == nil {
				if got != nil {
					t.Errorf("EffectiveOn() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("EffectiveOn() = nil, want %v", tt.want)
			}
			if (got.GitHubIssues != nil) != (tt.want.GitHubIssues != nil) {
				t.Errorf("EffectiveOn().GitHubIssues = %v, want %v", got.GitHubIssues, tt.want.GitHubIssues)
			}
			if (got.Cron != nil) != (tt.want.Cron != nil) {
				t.Errorf("EffectiveOn().Cron = %v, want %v", got.Cron, tt.want.Cron)
			}
			if got.Cron != nil && tt.want.Cron != nil && got.Cron.Schedule != tt.want.Cron.Schedule {
				t.Errorf("EffectiveOn().Cron.Schedule = %q, want %q", got.Cron.Schedule, tt.want.Cron.Schedule)
			}
		})
	}
}
