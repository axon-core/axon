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
	tasks := []axonv1alpha1.Task{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-one",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type: "claude-code",
			},
			Status: axonv1alpha1.TaskStatus{
				Phase: axonv1alpha1.TaskPhaseRunning,
			},
		},
	}

	var buf bytes.Buffer
	printTaskTable(&buf, tasks, false)
	output := buf.String()

	if !strings.Contains(output, "NAME") {
		t.Errorf("expected header NAME in output, got %q", output)
	}
	if strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected no NAMESPACE header when allNamespaces is false, got %q", output)
	}
	if !strings.Contains(output, "task-one") {
		t.Errorf("expected task-one in output, got %q", output)
	}
}

func TestPrintTaskTableAllNamespaces(t *testing.T) {
	tasks := []axonv1alpha1.Task{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-one",
				Namespace:         "ns-a",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type: "claude-code",
			},
			Status: axonv1alpha1.TaskStatus{
				Phase: axonv1alpha1.TaskPhaseRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "task-two",
				Namespace:         "ns-b",
				CreationTimestamp: metav1.NewTime(time.Now().Add(-2 * time.Hour)),
			},
			Spec: axonv1alpha1.TaskSpec{
				Type: "codex",
			},
			Status: axonv1alpha1.TaskStatus{
				Phase: axonv1alpha1.TaskPhaseSucceeded,
			},
		},
	}

	var buf bytes.Buffer
	printTaskTable(&buf, tasks, true)
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
