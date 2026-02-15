package controller

import (
	"context"
	"testing"

	axonv1alpha1 "github.com/axon-core/axon/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestParseGitHubOwnerRepo(t *testing.T) {
	tests := []struct {
		name      string
		repoURL   string
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "HTTPS URL",
			repoURL:   "https://github.com/axon-core/axon.git",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "HTTPS URL without .git",
			repoURL:   "https://github.com/axon-core/axon",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "HTTPS URL with trailing slash",
			repoURL:   "https://github.com/axon-core/axon/",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "SSH URL",
			repoURL:   "git@github.com:axon-core/axon.git",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "SSH URL without .git",
			repoURL:   "git@github.com:axon-core/axon",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "HTTPS URL with org",
			repoURL:   "https://github.com/my-org/my-repo.git",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
		{
			name:      "GitHub Enterprise HTTPS URL",
			repoURL:   "https://github.example.com/my-org/my-repo.git",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
		{
			name:      "GitHub Enterprise SSH URL",
			repoURL:   "git@github.example.com:my-org/my-repo.git",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo := parseGitHubOwnerRepo(tt.repoURL)
			if owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}

func TestParseGitHubRepo(t *testing.T) {
	tests := []struct {
		name      string
		repoURL   string
		wantHost  string
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "github.com HTTPS",
			repoURL:   "https://github.com/axon-core/axon.git",
			wantHost:  "github.com",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "github.com SSH",
			repoURL:   "git@github.com:axon-core/axon.git",
			wantHost:  "github.com",
			wantOwner: "axon-core",
			wantRepo:  "axon",
		},
		{
			name:      "GitHub Enterprise HTTPS",
			repoURL:   "https://github.example.com/my-org/my-repo.git",
			wantHost:  "github.example.com",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
		{
			name:      "GitHub Enterprise SSH",
			repoURL:   "git@github.example.com:my-org/my-repo.git",
			wantHost:  "github.example.com",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
		{
			name:      "GitHub Enterprise HTTPS without .git",
			repoURL:   "https://github.example.com/my-org/my-repo",
			wantHost:  "github.example.com",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
		{
			name:      "GitHub Enterprise with port",
			repoURL:   "https://github.example.com:8443/my-org/my-repo.git",
			wantHost:  "github.example.com:8443",
			wantOwner: "my-org",
			wantRepo:  "my-repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, owner, repo := parseGitHubRepo(tt.repoURL)
			if host != tt.wantHost {
				t.Errorf("host = %q, want %q", host, tt.wantHost)
			}
			if owner != tt.wantOwner {
				t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}

func TestGitHubAPIBaseURL(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{
			name: "empty host returns empty",
			host: "",
			want: "",
		},
		{
			name: "github.com returns empty",
			host: "github.com",
			want: "",
		},
		{
			name: "enterprise host",
			host: "github.example.com",
			want: "https://github.example.com/api/v3",
		},
		{
			name: "enterprise host with port",
			host: "github.example.com:8443",
			want: "https://github.example.com:8443/api/v3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gitHubAPIBaseURL(tt.host)
			if got != tt.want {
				t.Errorf("gitHubAPIBaseURL(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}

func TestBuildDeploymentWithEnterpriseURL(t *testing.T) {
	builder := NewDeploymentBuilder()

	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
		},
	}

	tests := []struct {
		name              string
		repoURL           string
		wantAPIBaseURLArg string
	}{
		{
			name:              "github.com repo does not include api-base-url arg",
			repoURL:           "https://github.com/axon-core/axon.git",
			wantAPIBaseURLArg: "",
		},
		{
			name:              "enterprise repo includes api-base-url arg",
			repoURL:           "https://github.example.com/my-org/my-repo.git",
			wantAPIBaseURLArg: "--github-api-base-url=https://github.example.com/api/v3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := &axonv1alpha1.WorkspaceSpec{
				Repo: tt.repoURL,
			}
			dep := builder.Build(ts, workspace, false)
			args := dep.Spec.Template.Spec.Containers[0].Args

			found := ""
			for _, arg := range args {
				if len(arg) > len("--github-api-base-url=") && arg[:len("--github-api-base-url=")] == "--github-api-base-url=" {
					found = arg
				}
			}

			if tt.wantAPIBaseURLArg == "" {
				if found != "" {
					t.Errorf("Expected no --github-api-base-url arg, got %q", found)
				}
			} else {
				if found != tt.wantAPIBaseURLArg {
					t.Errorf("Got arg %q, want %q", found, tt.wantAPIBaseURLArg)
				}
			}
		})
	}
}

func TestDeploymentBuilder_GitHubApp(t *testing.T) {
	builder := NewDeploymentBuilder()
	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
			TaskTemplate: axonv1alpha1.TaskTemplate{
				Type:         "claude-code",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{Name: "ws"},
			},
		},
	}
	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/axon-core/axon.git",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-app-creds",
		},
	}

	deploy := builder.Build(ts, workspace, true)

	if len(deploy.Spec.Template.Spec.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(deploy.Spec.Template.Spec.Containers))
	}

	if len(deploy.Spec.Template.Spec.InitContainers) != 1 {
		t.Fatalf("expected 1 init container, got %d", len(deploy.Spec.Template.Spec.InitContainers))
	}

	spawner := deploy.Spec.Template.Spec.Containers[0]
	refresher := deploy.Spec.Template.Spec.InitContainers[0]

	if spawner.Name != "spawner" {
		t.Errorf("container name = %q, want %q", spawner.Name, "spawner")
	}
	if refresher.Name != "token-refresher" {
		t.Errorf("init container name = %q, want %q", refresher.Name, "token-refresher")
	}

	if refresher.RestartPolicy == nil || *refresher.RestartPolicy != corev1.ContainerRestartPolicyAlways {
		t.Errorf("token-refresher RestartPolicy = %v, want %q", refresher.RestartPolicy, corev1.ContainerRestartPolicyAlways)
	}

	found := false
	for _, arg := range spawner.Args {
		if arg == "--github-token-file=/shared/token/GITHUB_TOKEN" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("spawner args missing --github-token-file flag: %v", spawner.Args)
	}

	for _, env := range spawner.Env {
		if env.Name == "GITHUB_TOKEN" {
			t.Error("spawner should not have GITHUB_TOKEN env var in GitHub App mode")
		}
	}

	if len(deploy.Spec.Template.Spec.Volumes) != 2 {
		t.Fatalf("expected 2 volumes, got %d", len(deploy.Spec.Template.Spec.Volumes))
	}

	if len(refresher.Env) != 2 {
		t.Fatalf("token-refresher expected 2 env vars, got %d", len(refresher.Env))
	}
	if refresher.Env[0].Name != "APP_ID" {
		t.Errorf("first env var = %q, want %q", refresher.Env[0].Name, "APP_ID")
	}
	if refresher.Env[1].Name != "INSTALLATION_ID" {
		t.Errorf("second env var = %q, want %q", refresher.Env[1].Name, "INSTALLATION_ID")
	}
}

func TestDeploymentBuilder_PAT(t *testing.T) {
	builder := NewDeploymentBuilder()
	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
			TaskTemplate: axonv1alpha1.TaskTemplate{
				Type:         "claude-code",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{Name: "ws"},
			},
		},
	}
	workspace := &axonv1alpha1.WorkspaceSpec{
		Repo: "https://github.com/axon-core/axon.git",
		SecretRef: &axonv1alpha1.SecretReference{
			Name: "github-token",
		},
	}

	deploy := builder.Build(ts, workspace, false)

	if len(deploy.Spec.Template.Spec.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(deploy.Spec.Template.Spec.Containers))
	}

	if len(deploy.Spec.Template.Spec.InitContainers) != 0 {
		t.Errorf("expected 0 init containers, got %d", len(deploy.Spec.Template.Spec.InitContainers))
	}

	spawner := deploy.Spec.Template.Spec.Containers[0]

	if len(spawner.Env) != 1 {
		t.Fatalf("expected 1 env var, got %d", len(spawner.Env))
	}
	if spawner.Env[0].Name != "GITHUB_TOKEN" {
		t.Errorf("env var name = %q, want %q", spawner.Env[0].Name, "GITHUB_TOKEN")
	}

	if len(deploy.Spec.Template.Spec.Volumes) != 0 {
		t.Errorf("expected 0 volumes, got %d", len(deploy.Spec.Template.Spec.Volumes))
	}
}

func boolPtr(v bool) *bool { return &v }

func TestUpdateDeployment_SuspendScalesDown(t *testing.T) {
	builder := NewDeploymentBuilder()
	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
			TaskTemplate: axonv1alpha1.TaskTemplate{
				Type:         "claude-code",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{Name: "ws"},
			},
			Suspend: boolPtr(true),
		},
	}

	// Build a deployment with replicas=1 (running state)
	deploy := builder.Build(ts, nil, false)
	if deploy.Spec.Replicas == nil || *deploy.Spec.Replicas != 1 {
		t.Fatalf("expected initial Replicas=1, got %v", deploy.Spec.Replicas)
	}

	// Create a reconciler with a fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(axonv1alpha1.AddToScheme(scheme))

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(ts, deploy).
		WithStatusSubresource(ts).
		Build()

	r := &TaskSpawnerReconciler{
		Client:            cl,
		Scheme:            scheme,
		DeploymentBuilder: builder,
	}

	// Call updateDeployment with desiredReplicas=0 (suspended)
	ctx := context.Background()
	if err := r.updateDeployment(ctx, ts, deploy, nil, false, 0); err != nil {
		t.Fatalf("updateDeployment error: %v", err)
	}

	// Verify the deployment was updated to 0 replicas
	var updated appsv1.Deployment
	if err := cl.Get(ctx, client.ObjectKeyFromObject(deploy), &updated); err != nil {
		t.Fatalf("getting deployment: %v", err)
	}
	if updated.Spec.Replicas == nil || *updated.Spec.Replicas != 0 {
		t.Errorf("expected Replicas=0 after suspend, got %v", updated.Spec.Replicas)
	}
}

func TestUpdateDeployment_ResumeScalesUp(t *testing.T) {
	builder := NewDeploymentBuilder()
	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
			TaskTemplate: axonv1alpha1.TaskTemplate{
				Type:         "claude-code",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{Name: "ws"},
			},
			Suspend: boolPtr(false),
		},
	}

	// Build a deployment with replicas=0 (suspended state)
	deploy := builder.Build(ts, nil, false)
	zero := int32(0)
	deploy.Spec.Replicas = &zero

	// Create a reconciler with a fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(axonv1alpha1.AddToScheme(scheme))

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(ts, deploy).
		WithStatusSubresource(ts).
		Build()

	r := &TaskSpawnerReconciler{
		Client:            cl,
		Scheme:            scheme,
		DeploymentBuilder: builder,
	}

	// Call updateDeployment with desiredReplicas=1 (resumed)
	ctx := context.Background()
	if err := r.updateDeployment(ctx, ts, deploy, nil, false, 1); err != nil {
		t.Fatalf("updateDeployment error: %v", err)
	}

	// Verify the deployment was updated to 1 replica
	var updated appsv1.Deployment
	if err := cl.Get(ctx, client.ObjectKeyFromObject(deploy), &updated); err != nil {
		t.Fatalf("getting deployment: %v", err)
	}
	if updated.Spec.Replicas == nil || *updated.Spec.Replicas != 1 {
		t.Errorf("expected Replicas=1 after resume, got %v", updated.Spec.Replicas)
	}
}

func TestUpdateDeployment_NoUpdateWhenReplicasMatch(t *testing.T) {
	builder := NewDeploymentBuilder()
	ts := &axonv1alpha1.TaskSpawner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-spawner",
			Namespace: "default",
		},
		Spec: axonv1alpha1.TaskSpawnerSpec{
			When: axonv1alpha1.When{
				GitHubIssues: &axonv1alpha1.GitHubIssues{},
			},
			TaskTemplate: axonv1alpha1.TaskTemplate{
				Type:         "claude-code",
				WorkspaceRef: &axonv1alpha1.WorkspaceReference{Name: "ws"},
			},
		},
	}

	// Build a deployment with replicas=1
	deploy := builder.Build(ts, nil, false)

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(axonv1alpha1.AddToScheme(scheme))

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(ts, deploy).
		WithStatusSubresource(ts).
		Build()

	r := &TaskSpawnerReconciler{
		Client:            cl,
		Scheme:            scheme,
		DeploymentBuilder: builder,
	}

	// Call updateDeployment with desiredReplicas=1 (no change needed)
	ctx := context.Background()
	if err := r.updateDeployment(ctx, ts, deploy, nil, false, 1); err != nil {
		t.Fatalf("updateDeployment error: %v", err)
	}

	// Verify the deployment still has 1 replica (no unnecessary update)
	var updated appsv1.Deployment
	if err := cl.Get(ctx, client.ObjectKeyFromObject(deploy), &updated); err != nil {
		t.Fatalf("getting deployment: %v", err)
	}
	if updated.Spec.Replicas == nil || *updated.Spec.Replicas != 1 {
		t.Errorf("expected Replicas=1 (unchanged), got %v", updated.Spec.Replicas)
	}
}
