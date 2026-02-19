package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	axonv1alpha1 "github.com/axon-core/axon/api/v1alpha1"
)

func TestPrintWorkspaceTable(t *testing.T) {
	workspaces := []axonv1alpha1.Workspace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-one",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/org/repo.git",
				Ref:  "main",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-two",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-24 * time.Hour)),
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/org/other.git",
			},
		},
	}

	var buf bytes.Buffer
	printWorkspaceTable(&buf, workspaces, false)
	output := buf.String()

	if !strings.Contains(output, "NAME") {
		t.Errorf("expected header NAME in output, got %q", output)
	}
	if !strings.Contains(output, "REPO") {
		t.Errorf("expected header REPO in output, got %q", output)
	}
	if strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected no NAMESPACE header when allNamespaces is false, got %q", output)
	}
	if !strings.Contains(output, "ws-one") {
		t.Errorf("expected ws-one in output, got %q", output)
	}
	if !strings.Contains(output, "ws-two") {
		t.Errorf("expected ws-two in output, got %q", output)
	}
	if !strings.Contains(output, "https://github.com/org/repo.git") {
		t.Errorf("expected repo URL in output, got %q", output)
	}
}

func TestPrintWorkspaceTableAllNamespaces(t *testing.T) {
	workspaces := []axonv1alpha1.Workspace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-one",
				Namespace:         "ns-a",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/org/repo.git",
				Ref:  "main",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-two",
				Namespace:         "ns-b",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-24 * time.Hour)),
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/org/other.git",
			},
		},
	}

	var buf bytes.Buffer
	printWorkspaceTable(&buf, workspaces, true)
	output := buf.String()

	if !strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected NAMESPACE header when allNamespaces is true, got %q", output)
	}
	if !strings.Contains(output, "ns-a") {
		t.Errorf("expected namespace ns-a in output, got %q", output)
	}
	if !strings.Contains(output, "ns-b") {
		t.Errorf("expected namespace ns-b in output, got %q", output)
	}
}

func TestPrintTaskTable(t *testing.T) {
	now := time.Now()
	startTime := metav1.NewTime(now.Add(-30 * time.Minute))
	tasks := []axonv1alpha1.Task{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-one",
				CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type:   "claude-code",
				Model:  "claude-sonnet-4-20250514",
				Branch: "feature/test",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{
					Name: "my-ws",
				},
				AgentConfigRef: &axonv1alpha1.AgentConfigReference{
					Name: "my-config",
				},
			},
			Status: axonv1alpha1.TaskStatus{
				Phase:     axonv1alpha1.TaskPhaseRunning,
				StartTime: &startTime,
			},
		},
	}

	var buf bytes.Buffer
	printTaskTable(&buf, tasks, false)
	output := buf.String()

	for _, header := range []string{"NAME", "TYPE", "PHASE", "BRANCH", "WORKSPACE", "AGENT CONFIG", "DURATION", "AGE"} {
		if !strings.Contains(output, header) {
			t.Errorf("expected header %s in output, got %q", header, output)
		}
	}
	if strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected no NAMESPACE header when allNamespaces is false, got %q", output)
	}
	for _, val := range []string{"task-one", "feature/test", "my-ws", "my-config"} {
		if !strings.Contains(output, val) {
			t.Errorf("expected %s in output, got %q", val, output)
		}
	}
}

func TestPrintTaskTableAllNamespaces(t *testing.T) {
	now := time.Now()
	startTime := metav1.NewTime(now.Add(-90 * time.Minute))
	completionTime := metav1.NewTime(now.Add(-60 * time.Minute))
	tasks := []axonv1alpha1.Task{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-one",
				Namespace:         "ns-a",
				CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type:   "claude-code",
				Branch: "feat/one",
			},
			Status: axonv1alpha1.TaskStatus{
				Phase: axonv1alpha1.TaskPhaseRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-two",
				Namespace:         "ns-b",
				CreationTimestamp: metav1.NewTime(now.Add(-2 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type: "codex",
			},
			Status: axonv1alpha1.TaskStatus{
				Phase:          axonv1alpha1.TaskPhaseSucceeded,
				StartTime:      &startTime,
				CompletionTime: &completionTime,
			},
		},
	}

	var buf bytes.Buffer
	printTaskTable(&buf, tasks, true)
	output := buf.String()

	for _, header := range []string{"NAMESPACE", "NAME", "TYPE", "PHASE", "BRANCH", "WORKSPACE", "AGENT CONFIG", "DURATION", "AGE"} {
		if !strings.Contains(output, header) {
			t.Errorf("expected header %s in output, got %q", header, output)
		}
	}
	for _, val := range []string{"ns-a", "ns-b", "feat/one"} {
		if !strings.Contains(output, val) {
			t.Errorf("expected %s in output, got %q", val, output)
		}
	}
}

func TestPrintTaskSpawnerTable(t *testing.T) {
	spawners := []axonv1alpha1.TaskSpawner{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "spawner-one",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpawnerSpec{
				When: axonv1alpha1.When{
					Cron: &axonv1alpha1.Cron{
						Schedule: "*/5 * * * *",
					},
				},
			},
			Status: axonv1alpha1.TaskSpawnerStatus{
				Phase: axonv1alpha1.TaskSpawnerPhaseRunning,
			},
		},
	}

	var buf bytes.Buffer
	printTaskSpawnerTable(&buf, spawners, false)
	output := buf.String()

	if !strings.Contains(output, "NAME") {
		t.Errorf("expected header NAME in output, got %q", output)
	}
	if strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected no NAMESPACE header when allNamespaces is false, got %q", output)
	}
	if !strings.Contains(output, "spawner-one") {
		t.Errorf("expected spawner-one in output, got %q", output)
	}
}

func TestPrintTaskSpawnerTableAllNamespaces(t *testing.T) {
	spawners := []axonv1alpha1.TaskSpawner{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "spawner-one",
				Namespace:         "ns-a",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpawnerSpec{
				When: axonv1alpha1.When{
					Cron: &axonv1alpha1.Cron{
						Schedule: "*/5 * * * *",
					},
				},
			},
			Status: axonv1alpha1.TaskSpawnerStatus{
				Phase: axonv1alpha1.TaskSpawnerPhaseRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "spawner-two",
				Namespace:         "ns-b",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpawnerSpec{
				When: axonv1alpha1.When{
					GitHubIssues: &axonv1alpha1.GitHubIssues{},
				},
				TaskTemplate: axonv1alpha1.TaskTemplate{
					WorkspaceRef: &axonv1alpha1.WorkspaceReference{
						Name: "my-ws",
					},
				},
			},
			Status: axonv1alpha1.TaskSpawnerStatus{
				Phase: axonv1alpha1.TaskSpawnerPhaseRunning,
			},
		},
	}

	var buf bytes.Buffer
	printTaskSpawnerTable(&buf, spawners, true)
	output := buf.String()

	if !strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected NAMESPACE header when allNamespaces is true, got %q", output)
	}
	if !strings.Contains(output, "ns-a") {
		t.Errorf("expected namespace ns-a in output, got %q", output)
	}
	if !strings.Contains(output, "ns-b") {
		t.Errorf("expected namespace ns-b in output, got %q", output)
	}
}

func TestPrintWorkspaceDetail(t *testing.T) {
	ws := &axonv1alpha1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-workspace",
			Namespace: "default",
		},
		Spec: axonv1alpha1.WorkspaceSpec{
			Repo: "https://github.com/org/repo.git",
			Ref:  "main",
			SecretRef: &axonv1alpha1.SecretReference{
				Name: "gh-token",
			},
		},
	}

	var buf bytes.Buffer
	printWorkspaceDetail(&buf, ws)
	output := buf.String()

	if !strings.Contains(output, "my-workspace") {
		t.Errorf("expected workspace name in output, got %q", output)
	}
	if !strings.Contains(output, "https://github.com/org/repo.git") {
		t.Errorf("expected repo URL in output, got %q", output)
	}
	if !strings.Contains(output, "main") {
		t.Errorf("expected ref in output, got %q", output)
	}
	if !strings.Contains(output, "gh-token") {
		t.Errorf("expected secret name in output, got %q", output)
	}
}

func TestPrintWorkspaceDetailWithoutOptionalFields(t *testing.T) {
	ws := &axonv1alpha1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minimal-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.WorkspaceSpec{
			Repo: "https://github.com/org/repo.git",
		},
	}

	var buf bytes.Buffer
	printWorkspaceDetail(&buf, ws)
	output := buf.String()

	if !strings.Contains(output, "minimal-ws") {
		t.Errorf("expected workspace name in output, got %q", output)
	}
	if strings.Contains(output, "Ref:") {
		t.Errorf("expected no Ref field when ref is empty, got %q", output)
	}
	if strings.Contains(output, "Secret:") {
		t.Errorf("expected no Secret field when secretRef is nil, got %q", output)
	}
}

func TestTaskDuration(t *testing.T) {
	now := time.Now()
	startTime := metav1.NewTime(now.Add(-30 * time.Minute))
	completionTime := metav1.NewTime(now.Add(-10 * time.Minute))

	tests := []struct {
		name   string
		status axonv1alpha1.TaskStatus
		want   string
	}{
		{
			name:   "no start time",
			status: axonv1alpha1.TaskStatus{},
			want:   "-",
		},
		{
			name: "completed task",
			status: axonv1alpha1.TaskStatus{
				StartTime:      &startTime,
				CompletionTime: &completionTime,
			},
			want: "20m",
		},
		{
			name: "running task",
			status: axonv1alpha1.TaskStatus{
				StartTime: &startTime,
			},
			want: "30m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskDuration(&tt.status)
			if got != tt.want {
				t.Errorf("taskDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintTaskDetail(t *testing.T) {
	now := time.Now()
	startTime := metav1.NewTime(now.Add(-30 * time.Minute))
	completionTime := metav1.NewTime(now.Add(-10 * time.Minute))
	ttl := int32(3600)
	timeout := int64(7200)

	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "full-task",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   "claude-code",
			Prompt: "Fix the bug",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "my-secret"},
			},
			Model:     "claude-sonnet-4-20250514",
			Image:     "custom-image:latest",
			Branch:    "feature/fix",
			DependsOn: []string{"task-a", "task-b"},
			WorkspaceRef: &axonv1alpha1.WorkspaceReference{
				Name: "my-ws",
			},
			AgentConfigRef: &axonv1alpha1.AgentConfigReference{
				Name: "my-config",
			},
			TTLSecondsAfterFinished: &ttl,
			PodOverrides: &axonv1alpha1.PodOverrides{
				ActiveDeadlineSeconds: &timeout,
			},
		},
		Status: axonv1alpha1.TaskStatus{
			Phase:          axonv1alpha1.TaskPhaseSucceeded,
			JobName:        "full-task-job",
			PodName:        "full-task-pod",
			StartTime:      &startTime,
			CompletionTime: &completionTime,
			Message:        "Task completed successfully",
			Outputs:        []string{"https://github.com/org/repo/pull/1"},
			Results:        map[string]string{"pr": "1"},
		},
	}

	var buf bytes.Buffer
	printTaskDetail(&buf, task)
	output := buf.String()

	for _, expected := range []string{
		"full-task",
		"claude-code",
		"Succeeded",
		"Fix the bug",
		"my-secret",
		"claude-sonnet-4-20250514",
		"custom-image:latest",
		"feature/fix",
		"task-a, task-b",
		"my-ws",
		"my-config",
		"3600s",
		"7200s",
		"full-task-job",
		"full-task-pod",
		"Duration:",
		"Task completed successfully",
		"https://github.com/org/repo/pull/1",
		"pr=1",
	} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q in output, got:\n%s", expected, output)
		}
	}
}

func TestPrintTaskDetailMinimal(t *testing.T) {
	task := &axonv1alpha1.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minimal-task",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpec{
			Type:   "claude-code",
			Prompt: "Do something",
			Credentials: axonv1alpha1.Credentials{
				Type:      axonv1alpha1.CredentialTypeAPIKey,
				SecretRef: axonv1alpha1.SecretReference{Name: "secret"},
			},
		},
		Status: axonv1alpha1.TaskStatus{
			Phase: axonv1alpha1.TaskPhasePending,
		},
	}

	var buf bytes.Buffer
	printTaskDetail(&buf, task)
	output := buf.String()

	if !strings.Contains(output, "minimal-task") {
		t.Errorf("expected task name in output, got:\n%s", output)
	}

	for _, absent := range []string{
		"Model:",
		"Image:",
		"Branch:",
		"Depends On:",
		"Workspace:",
		"Agent Config:",
		"TTL:",
		"Timeout:",
		"Job:",
		"Pod:",
		"Start Time:",
		"Completion Time:",
		"Duration:",
		"Message:",
		"Outputs:",
		"Results:",
	} {
		if strings.Contains(output, absent) {
			t.Errorf("expected no %s field for minimal task, got:\n%s", absent, output)
		}
	}
}
