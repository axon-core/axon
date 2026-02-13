package e2e

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/axon-core/axon/test/e2e/framework"
)

var _ = Describe("TaskSpawner", func() {
	f := framework.NewFramework("spawner")

	BeforeEach(func() {
		if githubToken == "" {
			Skip("GITHUB_TOKEN not set, skipping TaskSpawner e2e tests")
		}
	})

	// This test requires at least one open GitHub issue in axon-core/axon
	// with the "do-not-remove/e2e-anchor" label. See issue #117.
	It("should create a spawner Deployment and discover issues", func() {
		By("creating GitHub token secret")
		f.CreateSecret("github-token",
			"GITHUB_TOKEN="+githubToken)

		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Workspace resource with secretRef")
		wsYAML := `apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: e2e-spawner-workspace
spec:
  repo: https://github.com/axon-core/axon.git
  ref: main
  secretRef:
    name: github-token
`
		f.ApplyYAML(wsYAML)

		By("creating a TaskSpawner")
		tsYAML := `apiVersion: axon.io/v1alpha1
kind: TaskSpawner
metadata:
  name: spawner
spec:
  when:
    githubIssues:
      labels: [do-not-remove/e2e-anchor]
      excludeLabels: [e2e-exclude-placeholder]
      state: open
  taskTemplate:
    type: claude-code
    workspaceRef:
      name: e2e-spawner-workspace
    credentials:
      type: oauth
      secretRef:
        name: claude-credentials
    promptTemplate: "Fix: {{.Title}}\n{{.Body}}"
  pollInterval: 1m
`
		f.ApplyYAML(tsYAML)

		By("waiting for Deployment to become available")
		f.WaitForDeploymentAvailable("spawner")

		By("waiting for TaskSpawner phase to become Running")
		Eventually(func() string {
			return f.KubectlOutput("get", "taskspawner", "spawner", "-o", "jsonpath={.status.phase}")
		}, 3*time.Minute, 10*time.Second).Should(Equal("Running"))

		By("verifying at least one Task was created")
		Eventually(func() string {
			return f.KubectlOutput("get", "tasks", "-l", "axon.io/taskspawner=spawner", "-o", "name")
		}, 3*time.Minute, 10*time.Second).ShouldNot(BeEmpty())
	})

	It("should be accessible via CLI", func() {
		By("creating a Workspace resource")
		wsYAML := `apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: e2e-spawner-workspace
spec:
  repo: https://github.com/axon-core/axon.git
`
		f.ApplyYAML(wsYAML)

		By("creating a TaskSpawner")
		tsYAML := `apiVersion: axon.io/v1alpha1
kind: TaskSpawner
metadata:
  name: spawner
spec:
  when:
    githubIssues: {}
  taskTemplate:
    type: claude-code
    workspaceRef:
      name: e2e-spawner-workspace
    credentials:
      type: oauth
      secretRef:
        name: claude-credentials
  pollInterval: 5m
`
		f.ApplyYAML(tsYAML)

		By("verifying axon get taskspawners lists it")
		output := framework.AxonOutput("get", "taskspawners", "-n", f.Namespace)
		Expect(output).To(ContainSubstring("spawner"))

		By("verifying axon get taskspawner shows detail")
		output = framework.AxonOutput("get", "taskspawner", "spawner", "-n", f.Namespace)
		Expect(output).To(ContainSubstring("spawner"))
		Expect(output).To(ContainSubstring("GitHub Issues"))

		By("verifying YAML output for a single taskspawner")
		output = framework.AxonOutput("get", "taskspawner", "spawner", "-n", f.Namespace, "-o", "yaml")
		Expect(output).To(ContainSubstring("apiVersion: axon.io/v1alpha1"))
		Expect(output).To(ContainSubstring("kind: TaskSpawner"))
		Expect(output).To(ContainSubstring("name: spawner"))

		By("verifying JSON output for a single taskspawner")
		output = framework.AxonOutput("get", "taskspawner", "spawner", "-n", f.Namespace, "-o", "json")
		Expect(output).To(ContainSubstring(`"apiVersion": "axon.io/v1alpha1"`))
		Expect(output).To(ContainSubstring(`"kind": "TaskSpawner"`))
		Expect(output).To(ContainSubstring(`"name": "spawner"`))

		By("deleting via kubectl")
		f.KubectlInNs("delete", "taskspawner", "spawner")

		By("verifying it disappears from list")
		Eventually(func() string {
			return framework.AxonOutput("get", "taskspawners", "-n", f.Namespace)
		}, 30*time.Second, time.Second).ShouldNot(ContainSubstring("spawner"))
	})
})

var _ = Describe("Cron TaskSpawner", func() {
	f := framework.NewFramework("cron")

	It("should create a spawner Deployment and discover cron ticks", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a cron TaskSpawner with every-minute schedule")
		tsYAML := `apiVersion: axon.io/v1alpha1
kind: TaskSpawner
metadata:
  name: cron-spawner
spec:
  when:
    cron:
      schedule: "* * * * *"
  taskTemplate:
    type: claude-code
    model: ` + testModel + `
    credentials:
      type: oauth
      secretRef:
        name: claude-credentials
    promptTemplate: "Cron triggered at {{.Time}} (schedule: {{.Schedule}}). Print 'Hello from cron'"
  pollInterval: 1m
`
		f.ApplyYAML(tsYAML)

		By("waiting for Deployment to become available")
		f.WaitForDeploymentAvailable("cron-spawner")

		By("waiting for TaskSpawner phase to become Running")
		Eventually(func() string {
			return f.KubectlOutput("get", "taskspawner", "cron-spawner", "-o", "jsonpath={.status.phase}")
		}, 3*time.Minute, 10*time.Second).Should(Equal("Running"))

		By("verifying at least one Task was created")
		Eventually(func() string {
			return f.KubectlOutput("get", "tasks", "-l", "axon.io/taskspawner=cron-spawner", "-o", "name")
		}, 3*time.Minute, 10*time.Second).ShouldNot(BeEmpty())
	})

	It("should be accessible via CLI with cron source info", func() {
		By("creating a cron TaskSpawner")
		tsYAML := `apiVersion: axon.io/v1alpha1
kind: TaskSpawner
metadata:
  name: cron-spawner
spec:
  when:
    cron:
      schedule: "0 9 * * 1"
  taskTemplate:
    type: claude-code
    credentials:
      type: oauth
      secretRef:
        name: claude-credentials
  pollInterval: 5m
`
		f.ApplyYAML(tsYAML)

		By("verifying axon get taskspawners lists it")
		output := framework.AxonOutput("get", "taskspawners", "-n", f.Namespace)
		Expect(output).To(ContainSubstring("cron-spawner"))

		By("verifying axon get taskspawner shows cron detail")
		output = framework.AxonOutput("get", "taskspawner", "cron-spawner", "-n", f.Namespace)
		Expect(output).To(ContainSubstring("cron-spawner"))
		Expect(output).To(ContainSubstring("Cron"))
		Expect(output).To(ContainSubstring("0 9 * * 1"))

		By("deleting via kubectl")
		f.KubectlInNs("delete", "taskspawner", "cron-spawner")

		By("verifying it disappears from list")
		Eventually(func() string {
			return framework.AxonOutput("get", "taskspawners", "-n", f.Namespace)
		}, 30*time.Second, time.Second).ShouldNot(ContainSubstring("cron-spawner"))
	})
})

var _ = Describe("get taskspawner", func() {
	It("should succeed with 'taskspawners' alias", func() {
		framework.AxonOutput("get", "taskspawners")
	})

	It("should succeed with 'ts' alias", func() {
		framework.AxonOutput("get", "ts")
	})

	It("should succeed with 'taskspawner' subcommand", func() {
		framework.AxonOutput("get", "taskspawner")
	})

	It("should fail for a nonexistent taskspawner", func() {
		framework.AxonFail("get", "taskspawner", "nonexistent-spawner")
	})

	It("should output taskspawner list in YAML format", func() {
		output := framework.AxonOutput("get", "taskspawners", "-o", "yaml")
		Expect(output).To(ContainSubstring("apiVersion: axon.io/v1alpha1"))
		Expect(output).To(ContainSubstring("kind: TaskSpawnerList"))
	})

	It("should output taskspawner list in JSON format", func() {
		output := framework.AxonOutput("get", "taskspawners", "-o", "json")
		Expect(output).To(ContainSubstring(`"apiVersion": "axon.io/v1alpha1"`))
		Expect(output).To(ContainSubstring(`"kind": "TaskSpawnerList"`))
	})

	It("should fail with unknown output format", func() {
		framework.AxonFail("get", "taskspawners", "-o", "invalid")
	})
})
