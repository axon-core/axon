package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

func TestPrintWorkspaceTable(t *testing.T) {
	workspaces := []axonv1alpha1.Workspace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-1",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-time.Hour)},
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/example/repo.git",
				Ref:  "main",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "ws-2",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(-24 * time.Hour)},
			},
			Spec: axonv1alpha1.WorkspaceSpec{
				Repo: "https://github.com/example/repo2.git",
			},
		},
	}

	var buf bytes.Buffer
	printWorkspaceTable(&buf, workspaces)
	output := buf.String()

	if !strings.Contains(output, "NAME") || !strings.Contains(output, "REPO") || !strings.Contains(output, "REF") || !strings.Contains(output, "AGE") {
		t.Fatalf("expected table headers, got: %s", output)
	}
	if !strings.Contains(output, "ws-1") {
		t.Fatalf("expected ws-1 in output, got: %s", output)
	}
	if !strings.Contains(output, "ws-2") {
		t.Fatalf("expected ws-2 in output, got: %s", output)
	}
	if !strings.Contains(output, "https://github.com/example/repo.git") {
		t.Fatalf("expected repo URL in output, got: %s", output)
	}
}

func TestPrintWorkspaceDetail(t *testing.T) {
	ws := &axonv1alpha1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.WorkspaceSpec{
			Repo: "https://github.com/example/repo.git",
			Ref:  "main",
			SecretRef: &axonv1alpha1.SecretReference{
				Name: "my-secret",
			},
		},
	}

	var buf bytes.Buffer
	printWorkspaceDetail(&buf, ws)
	output := buf.String()

	if !strings.Contains(output, "test-ws") {
		t.Fatalf("expected workspace name in output, got: %s", output)
	}
	if !strings.Contains(output, "default") {
		t.Fatalf("expected namespace in output, got: %s", output)
	}
	if !strings.Contains(output, "https://github.com/example/repo.git") {
		t.Fatalf("expected repo in output, got: %s", output)
	}
	if !strings.Contains(output, "main") {
		t.Fatalf("expected ref in output, got: %s", output)
	}
	if !strings.Contains(output, "my-secret") {
		t.Fatalf("expected secret in output, got: %s", output)
	}
}

func TestPrintWorkspaceDetail_NoOptionalFields(t *testing.T) {
	ws := &axonv1alpha1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ws",
			Namespace: "default",
		},
		Spec: axonv1alpha1.WorkspaceSpec{
			Repo: "https://github.com/example/repo.git",
		},
	}

	var buf bytes.Buffer
	printWorkspaceDetail(&buf, ws)
	output := buf.String()

	if strings.Contains(output, "Ref:") {
		t.Fatalf("unexpected Ref field in output for empty ref, got: %s", output)
	}
	if strings.Contains(output, "Secret:") {
		t.Fatalf("unexpected Secret field in output for nil secretRef, got: %s", output)
	}
}
