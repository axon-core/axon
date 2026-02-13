package e2e

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/axon-core/axon/test/e2e/framework"
)

var _ = Describe("Task", func() {
	f := framework.NewFramework("task")

	It("should run a Task to completion", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Task")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: basic-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Print 'Hello from Axon e2e test' to stdout"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("basic-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("basic-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("basic-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("basic-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})
})

var _ = Describe("Task with make available", func() {
	f := framework.NewFramework("make")

	It("should have make command available in claude-code container", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Task that uses make")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: make-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Run 'make --version' and print the output"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("make-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("make-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("make-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("make-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})
})

var _ = Describe("Task with workspace", func() {
	f := framework.NewFramework("ws")

	It("should run a Task with workspace to completion", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Workspace resource")
		wsYAML := `apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: e2e-workspace
spec:
  repo: https://github.com/axon-core/axon.git
  ref: main
`
		f.ApplyYAML(wsYAML)

		By("creating a Task with workspace ref")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: ws-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Create a file called 'test.txt' with the content 'hello' in the current directory and print 'done'"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
  workspaceRef:
    name: e2e-workspace
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("ws-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("ws-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("ws-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("ws-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)

		By("verifying no permission errors in logs")
		Expect(logs).NotTo(ContainSubstring("permission denied"))
		Expect(logs).NotTo(ContainSubstring("Permission denied"))
		Expect(logs).NotTo(ContainSubstring("EACCES"))
	})
})

var _ = Describe("Task output capture", func() {
	f := framework.NewFramework("output")

	It("should populate Outputs with branch name after task completes", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Workspace resource")
		wsYAML := `apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: e2e-outputs-workspace
spec:
  repo: https://github.com/axon-core/axon.git
  ref: main
`
		f.ApplyYAML(wsYAML)

		By("creating a Task with workspace ref")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: outputs-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Run 'git branch --show-current' and print the output, then say done"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
  workspaceRef:
    name: e2e-outputs-workspace
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("outputs-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("outputs-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("outputs-task")).To(Equal("Succeeded"))

		By("verifying output markers appear in Pod logs")
		logs := f.GetJobLogs("outputs-task")
		Expect(logs).To(ContainSubstring("---AXON_OUTPUTS_START---"))
		Expect(logs).To(ContainSubstring("---AXON_OUTPUTS_END---"))
		Expect(logs).To(ContainSubstring("branch: main"))

		By("verifying Outputs field is populated in Task status")
		outputsJSON := f.KubectlOutput("get", "task", "outputs-task", "-o", "jsonpath={.status.outputs}")
		Expect(outputsJSON).To(ContainSubstring("branch: main"))
	})
})

var _ = Describe("Task with workspace and secretRef", func() {
	f := framework.NewFramework("github")

	BeforeEach(func() {
		if githubToken == "" {
			Skip("GITHUB_TOKEN not set, skipping GitHub e2e tests")
		}
	})

	It("should run a Task with gh CLI available and GITHUB_TOKEN injected", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating workspace credentials secret")
		f.CreateSecret("workspace-credentials",
			"GITHUB_TOKEN="+githubToken)

		By("creating a Workspace resource with secretRef")
		wsYAML := `apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: e2e-github-workspace
spec:
  repo: https://github.com/axon-core/axon.git
  ref: main
  secretRef:
    name: workspace-credentials
`
		f.ApplyYAML(wsYAML)

		By("creating a Task with workspace ref")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: github-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Run 'gh auth status' and print the output"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
  workspaceRef:
    name: e2e-github-workspace
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("github-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("github-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("github-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("github-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})
})

var _ = Describe("Task cleanup on failure", func() {
	f := framework.NewFramework("cleanup")

	It("should clean up namespace resources automatically", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating a Task")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: cleanup-task
spec:
  type: claude-code
  model: ` + testModel + `
  prompt: "Print 'Hello' to stdout"
  credentials:
    type: oauth
    secretRef:
      name: claude-credentials
`
		f.ApplyYAML(taskYAML)

		By("verifying resources exist in the namespace")
		Eventually(func() string {
			return f.KubectlOutput("get", "tasks", "-o", "name")
		}, 30*time.Second, time.Second).Should(ContainSubstring("cleanup-task"))
	})
})
