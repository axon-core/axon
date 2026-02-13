package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/axon-core/axon/test/e2e/framework"
)

// openCodeTestModel uses a free OpenCode model so e2e tests require no authentication.
const openCodeTestModel = "opencode/big-pickle"

var _ = Describe("OpenCode Task", func() {
	f := framework.NewFramework("opencode")

	It("should run an OpenCode Task to completion", func() {
		By("creating credentials secret (empty key for free OpenCode model)")
		f.CreateSecret("opencode-credentials",
			"OPENCODE_API_KEY=")

		By("creating an OpenCode Task")
		taskYAML := `apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: opencode-task
spec:
  type: opencode
  model: ` + openCodeTestModel + `
  prompt: "Print 'Hello from OpenCode e2e test' to stdout"
  credentials:
    type: api-key
    secretRef:
      name: opencode-credentials
`
		f.ApplyYAML(taskYAML)

		By("waiting for Job to be created")
		f.WaitForJobCreation("opencode-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("opencode-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("opencode-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("opencode-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})
})
