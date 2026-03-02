package e2e

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kelosv1alpha1 "github.com/kelos-dev/kelos/api/v1alpha1"
	"github.com/kelos-dev/kelos/test/e2e/framework"
)

var _ = Describe("Task with GitHub plugin", func() {
	f := framework.NewFramework("github-plugin")

	BeforeEach(func() {
		if oauthToken == "" {
			Skip("CLAUDE_CODE_OAUTH_TOKEN not set")
		}
	})

	It("should clone a public GitHub plugin and make it available to the agent", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		By("creating an AgentConfig with a GitHub plugin")
		f.CreateAgentConfig(&kelosv1alpha1.AgentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "github-plugin-config",
			},
			Spec: kelosv1alpha1.AgentConfigSpec{
				Plugins: []kelosv1alpha1.PluginSpec{
					{
						Name: "test-plugin",
						GitHub: &kelosv1alpha1.GitHubPluginSource{
							Repo: "kelos-dev/kelos",
						},
					},
				},
			},
		})

		By("verifying AgentConfig was created with GitHub plugin")
		ac, err := f.KelosClientset.ApiV1alpha1().AgentConfigs(f.Namespace).Get(
			context.TODO(), "github-plugin-config", metav1.GetOptions{},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(ac.Spec.Plugins).To(HaveLen(1))
		Expect(ac.Spec.Plugins[0].Name).To(Equal("test-plugin"))
		Expect(ac.Spec.Plugins[0].GitHub).NotTo(BeNil())
		Expect(ac.Spec.Plugins[0].GitHub.Repo).To(Equal("kelos-dev/kelos"))

		By("creating a Task that references the AgentConfig")
		f.CreateTask(&kelosv1alpha1.Task{
			ObjectMeta: metav1.ObjectMeta{
				Name: "github-plugin-task",
			},
			Spec: kelosv1alpha1.TaskSpec{
				Type:           "claude-code",
				Model:          testModel,
				Prompt:         "List the contents of /kelos/plugin/test-plugin/ directory and print them. Then print 'PLUGIN_DIR_EXISTS' if the directory exists.",
				AgentConfigRef: &kelosv1alpha1.AgentConfigReference{Name: "github-plugin-config"},
				Credentials: kelosv1alpha1.Credentials{
					Type:      kelosv1alpha1.CredentialTypeOAuth,
					SecretRef: kelosv1alpha1.SecretReference{Name: "claude-credentials"},
				},
			},
		})

		By("waiting for Job to be created")
		f.WaitForJobCreation("github-plugin-task")

		By("verifying Job has plugin-setup init container")
		job, err := f.Clientset.BatchV1().Jobs(f.Namespace).Get(
			context.TODO(), "github-plugin-task", metav1.GetOptions{},
		)
		Expect(err).NotTo(HaveOccurred())
		initContainerNames := make([]string, 0)
		for _, ic := range job.Spec.Template.Spec.InitContainers {
			initContainerNames = append(initContainerNames, ic.Name)
		}
		Expect(initContainerNames).To(ContainElement("plugin-setup"))

		By("waiting for Job to complete")
		f.WaitForJobCompletion("github-plugin-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("github-plugin-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("github-plugin-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})

	It("should clone a GitHub plugin with a specific ref", func() {
		By("creating OAuth credentials secret")
		f.CreateSecret("claude-credentials",
			"CLAUDE_CODE_OAUTH_TOKEN="+oauthToken)

		ref := "main"
		By(fmt.Sprintf("creating an AgentConfig with a GitHub plugin pinned to ref %q", ref))
		f.CreateAgentConfig(&kelosv1alpha1.AgentConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "github-plugin-ref-config",
			},
			Spec: kelosv1alpha1.AgentConfigSpec{
				Plugins: []kelosv1alpha1.PluginSpec{
					{
						Name: "ref-plugin",
						GitHub: &kelosv1alpha1.GitHubPluginSource{
							Repo: "kelos-dev/kelos",
							Ref:  &ref,
						},
					},
				},
			},
		})

		By("creating a Task that references the AgentConfig")
		f.CreateTask(&kelosv1alpha1.Task{
			ObjectMeta: metav1.ObjectMeta{
				Name: "github-plugin-ref-task",
			},
			Spec: kelosv1alpha1.TaskSpec{
				Type:           "claude-code",
				Model:          testModel,
				Prompt:         "List the contents of /kelos/plugin/ref-plugin/ directory. Print 'REF_PLUGIN_EXISTS' if the directory exists and has files.",
				AgentConfigRef: &kelosv1alpha1.AgentConfigReference{Name: "github-plugin-ref-config"},
				Credentials: kelosv1alpha1.Credentials{
					Type:      kelosv1alpha1.CredentialTypeOAuth,
					SecretRef: kelosv1alpha1.SecretReference{Name: "claude-credentials"},
				},
			},
		})

		By("waiting for Job to be created")
		f.WaitForJobCreation("github-plugin-ref-task")

		By("waiting for Job to complete")
		f.WaitForJobCompletion("github-plugin-ref-task")

		By("verifying Task status is Succeeded")
		Expect(f.GetTaskPhase("github-plugin-ref-task")).To(Equal("Succeeded"))

		By("getting Job logs")
		logs := f.GetJobLogs("github-plugin-ref-task")
		GinkgoWriter.Printf("Job logs:\n%s\n", logs)
	})
})

var _ = Describe("CLI GitHub plugin dry-run", func() {
	It("should generate correct YAML with --github-plugin flag", func() {
		output := framework.KelosOutput("create", "agentconfig", "test-gh-plugin",
			"--github-plugin", "my-plugin=acme/tools",
			"--dry-run",
		)

		Expect(output).To(ContainSubstring("kind: AgentConfig"))
		Expect(output).To(ContainSubstring("name: my-plugin"))
		Expect(output).To(ContainSubstring("repo: acme/tools"))
	})

	It("should generate correct YAML with --github-plugin flag including ref and host", func() {
		output := framework.KelosOutput("create", "agentconfig", "test-gh-plugin-ref",
			"--github-plugin", "my-plugin=acme/tools@v1.0,host=github.corp.com,secret=my-secret",
			"--dry-run",
		)

		Expect(output).To(ContainSubstring("kind: AgentConfig"))
		Expect(output).To(ContainSubstring("name: my-plugin"))
		Expect(output).To(ContainSubstring("repo: acme/tools"))
		Expect(output).To(ContainSubstring("ref: v1.0"))
		Expect(output).To(ContainSubstring("host: github.corp.com"))
		Expect(output).To(ContainSubstring("name: my-secret"))
	})

	It("should generate correct YAML combining --github-plugin with --skill", func() {
		output := framework.KelosOutput("create", "agentconfig", "test-combined",
			"--skill", "review=Review the code",
			"--github-plugin", "external=acme/tools",
			"--dry-run",
		)

		Expect(output).To(ContainSubstring("kind: AgentConfig"))
		Expect(output).To(ContainSubstring("name: kelos"))
		Expect(output).To(ContainSubstring("name: external"))
		Expect(output).To(ContainSubstring("repo: acme/tools"))
		Expect(output).To(ContainSubstring("name: review"))
	})

	It("should reject duplicate --github-plugin names", func() {
		framework.KelosFail("create", "agentconfig", "test-dup",
			"--github-plugin", "dup=acme/tools",
			"--github-plugin", "dup=acme/other",
			"--dry-run",
		)
	})

	It("should reject --github-plugin name colliding with inline kelos plugin", func() {
		framework.KelosFail("create", "agentconfig", "test-collision",
			"--skill", "review=Review the code",
			"--github-plugin", "kelos=acme/tools",
			"--dry-run",
		)
	})
})
