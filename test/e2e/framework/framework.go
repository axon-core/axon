package framework

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Framework provides a per-test namespace environment for e2e tests,
// similar to the Kubernetes e2e framework. Each Framework instance
// creates a unique namespace before each test and tears it down after.
type Framework struct {
	// BaseName is used as a prefix when generating the namespace name.
	BaseName string

	// Namespace is the Kubernetes namespace created for the current test.
	// It is set during BeforeEach and cleared during AfterEach.
	Namespace string
}

// NewFramework creates a new Framework and registers Ginkgo lifecycle
// hooks (BeforeEach/AfterEach) that manage the test namespace.
func NewFramework(baseName string) *Framework {
	f := &Framework{
		BaseName: baseName,
	}

	BeforeEach(f.beforeEach)
	AfterEach(f.afterEach)

	return f
}

func (f *Framework) beforeEach() {
	suffix := randomSuffix(6)
	f.Namespace = fmt.Sprintf("e2e-%s-%s", f.BaseName, suffix)
	// Kubernetes namespace names must be at most 63 characters and
	// conform to RFC 1123 DNS labels (no trailing hyphens).
	if len(f.Namespace) > 63 {
		f.Namespace = strings.TrimRight(f.Namespace[:63], "-")
	}

	By(fmt.Sprintf("Creating test namespace %s", f.Namespace))
	cmd := exec.Command("kubectl", "create", "namespace", f.Namespace)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	Expect(cmd.Run()).To(Succeed())
}

// randomSuffix returns a random hex string of the given length.
func randomSuffix(n int) string {
	b := make([]byte, (n+1)/2)
	_, err := rand.Read(b)
	if err != nil {
		// Fall back to timestamp if crypto/rand fails.
		s := fmt.Sprintf("%d", time.Now().UnixNano())
		if len(s) > n {
			s = s[:n]
		}
		return s
	}
	return fmt.Sprintf("%x", b)[:n]
}

func (f *Framework) afterEach() {
	if f.Namespace == "" {
		return
	}

	if CurrentSpecReport().Failed() {
		By("Collecting debug info on failure")
		f.Kubectl("get", "all", "-n", f.Namespace, "-o", "wide")
		f.Kubectl("get", "tasks.axon.io", "-n", f.Namespace, "-o", "yaml")
		f.Kubectl("get", "taskspawners.axon.io", "-n", f.Namespace, "-o", "yaml")
		f.Kubectl("logs", "-n", "axon-system", "deployment/axon-controller-manager", "--tail=50")
	}

	By(fmt.Sprintf("Deleting test namespace %s", f.Namespace))
	f.Kubectl("delete", "namespace", f.Namespace, "--ignore-not-found")

	f.Namespace = ""
}

// Kubectl executes a kubectl command with output directed to GinkgoWriter.
// It does NOT fail the test on error (fire-and-forget).
func (f *Framework) Kubectl(args ...string) {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	_ = cmd.Run()
}

// KubectlInNs executes a kubectl command scoped to the test namespace.
// It does NOT fail the test on error (fire-and-forget).
func (f *Framework) KubectlInNs(args ...string) {
	nsArgs := make([]string, 0, len(args)+2)
	nsArgs = append(nsArgs, "-n", f.Namespace)
	nsArgs = append(nsArgs, args...)
	f.Kubectl(nsArgs...)
}

// KubectlWithInput executes a kubectl command with optional stdin and
// returns an error. The command is scoped to the test namespace.
func (f *Framework) KubectlWithInput(input string, args ...string) error {
	nsArgs := make([]string, 0, len(args)+2)
	nsArgs = append(nsArgs, "-n", f.Namespace)
	nsArgs = append(nsArgs, args...)
	cmd := exec.Command("kubectl", nsArgs...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	return cmd.Run()
}

// KubectlOutput executes a kubectl command scoped to the test namespace
// and returns its stdout. It fails the test on error.
func (f *Framework) KubectlOutput(args ...string) string {
	nsArgs := make([]string, 0, len(args)+2)
	nsArgs = append(nsArgs, "-n", f.Namespace)
	nsArgs = append(nsArgs, args...)
	cmd := exec.Command("kubectl", nsArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	return strings.TrimSpace(out.String())
}

// CreateSecret creates a generic secret from literal key-value pairs
// in the test namespace.
func (f *Framework) CreateSecret(name string, literals ...string) {
	args := []string{"create", "secret", "generic", name}
	for _, l := range literals {
		args = append(args, "--from-literal="+l)
	}
	Expect(f.KubectlWithInput("", args...)).To(Succeed())
}

// ApplyYAML applies the given YAML manifest in the test namespace.
func (f *Framework) ApplyYAML(yaml string) {
	Expect(f.KubectlWithInput(yaml, "apply", "-f", "-")).To(Succeed())
}

// WaitForJobCreation waits for a Job with the given name to appear.
func (f *Framework) WaitForJobCreation(name string) {
	Eventually(func() error {
		return f.KubectlWithInput("", "get", "job", name)
	}, 30*time.Second, time.Second).Should(Succeed())
}

// WaitForJobCompletion waits for a Job to reach the complete condition.
func (f *Framework) WaitForJobCompletion(name string) {
	Eventually(func() error {
		return f.KubectlWithInput("", "wait", "--for=condition=complete", "job/"+name, "--timeout=10s")
	}, 5*time.Minute, 10*time.Second).Should(Succeed())
}

// WaitForDeploymentAvailable waits for a Deployment to reach the available condition.
func (f *Framework) WaitForDeploymentAvailable(name string) {
	Eventually(func() error {
		return f.KubectlWithInput("", "wait", "--for=condition=available", "deployment/"+name, "--timeout=10s")
	}, 2*time.Minute, 10*time.Second).Should(Succeed())
}

// GetTaskPhase returns the phase of a Task.
func (f *Framework) GetTaskPhase(name string) string {
	return f.KubectlOutput("get", "task", name, "-o", "jsonpath={.status.phase}")
}

// GetJobLogs returns the logs of a Job.
func (f *Framework) GetJobLogs(name string) string {
	return f.KubectlOutput("logs", "job/"+name)
}

// AxonBin returns the path to the axon binary.
func AxonBin() string {
	if bin := os.Getenv("AXON_BIN"); bin != "" {
		return bin
	}
	return "axon"
}

// Axon executes an axon CLI command with output directed to GinkgoWriter.
// It fails the test on error.
func Axon(args ...string) {
	cmd := exec.Command(AxonBin(), args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
}

// AxonOutput executes an axon CLI command and returns its stdout.
// It fails the test on error.
func AxonOutput(args ...string) string {
	cmd := exec.Command(AxonBin(), args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	return strings.TrimSpace(out.String())
}

// AxonOutputWithStderr executes an axon CLI command and returns both
// stdout and stderr.
func AxonOutputWithStderr(args ...string) (string, string) {
	cmd := exec.Command(AxonBin(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String())
}

// AxonFail executes an axon CLI command and expects it to fail.
func AxonFail(args ...string) {
	cmd := exec.Command(AxonBin(), args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	err := cmd.Run()
	Expect(err).To(HaveOccurred())
}

// AxonCommand creates an exec.Cmd for the axon binary without running it.
func AxonCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(AxonBin(), args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	return cmd
}
