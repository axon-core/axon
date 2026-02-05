package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

// NewCreateCmd creates a new create command
func NewCreateCmd() *cobra.Command {
	var (
		taskType   string
		prompt     string
		credType   string
		secretName string
		model      string
		taskName   string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task",
		Long:  `Create a new AI agent task in Kubernetes.`,
		Example: `  # Create a task with API key authentication
  axonctl create --type claude-code --prompt "Create a hello world in Python" \
    --cred-type api-key --secret anthropic-api-key

  # Create a task with OAuth authentication and custom model
  axonctl create --type claude-code --prompt "Fix the bug in main.go" \
    --cred-type oauth --secret claude-oauth --model claude-sonnet-4-20250514`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if taskType == "" || prompt == "" || credType == "" || secretName == "" {
				return fmt.Errorf("--type, --prompt, --cred-type, and --secret are required")
			}

			client, err := getAxonClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Generate task name if not provided
			if taskName == "" {
				taskName = fmt.Sprintf("task-%s", metav1.Now().Format("20060102150405"))
			}

			task := &axonv1alpha1.Task{
				ObjectMeta: metav1.ObjectMeta{
					Name:      taskName,
					Namespace: namespace,
				},
				Spec: axonv1alpha1.TaskSpec{
					Type:   taskType,
					Prompt: prompt,
					Credentials: axonv1alpha1.Credentials{
						Type: axonv1alpha1.CredentialType(credType),
						SecretRef: axonv1alpha1.SecretReference{
							Name: secretName,
						},
					},
					Model: model,
				},
			}

			ctx := context.Background()
			if err := client.Create(ctx, task); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("Task %s created successfully in namespace %s\n", taskName, namespace)
			return nil
		},
	}

	cmd.Flags().StringVar(&taskType, "type", "", "Agent type (e.g., claude-code)")
	cmd.Flags().StringVar(&prompt, "prompt", "", "Task prompt for the agent")
	cmd.Flags().StringVar(&credType, "cred-type", "", "Credential type (api-key or oauth)")
	cmd.Flags().StringVar(&secretName, "secret", "", "Name of the secret containing credentials")
	cmd.Flags().StringVar(&model, "model", "", "Optional model override")
	cmd.Flags().StringVar(&taskName, "name", "", "Task name (auto-generated if not provided)")

	return cmd
}
