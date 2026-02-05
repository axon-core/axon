package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

// NewGetCmd creates a new get command
func NewGetCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get TASK_NAME",
		Short: "Get task details",
		Long:  `Get detailed information about a specific task.`,
		Example: `  # Get task details
  axonctl get my-task

  # Get task details with YAML output
  axonctl get my-task -o yaml`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]

			client, err := getAxonClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			ctx := context.Background()
			task := &axonv1alpha1.Task{}
			if err := client.Get(ctx, types.NamespacedName{
				Name:      taskName,
				Namespace: namespace,
			}, task); err != nil {
				return fmt.Errorf("failed to get task: %w", err)
			}

			if outputFormat == "yaml" || outputFormat == "json" {
				// For YAML/JSON output format
				fmt.Printf("Name: %s\n", task.Name)
				fmt.Printf("Namespace: %s\n", task.Namespace)
				fmt.Printf("Type: %s\n", task.Spec.Type)
				fmt.Printf("Phase: %s\n", task.Status.Phase)
				if task.Spec.Model != "" {
					fmt.Printf("Model: %s\n", task.Spec.Model)
				}
				fmt.Printf("Prompt: %s\n", task.Spec.Prompt)
				if task.Status.JobName != "" {
					fmt.Printf("Job: %s\n", task.Status.JobName)
				}
				if task.Status.PodName != "" {
					fmt.Printf("Pod: %s\n", task.Status.PodName)
				}
				if task.Status.StartTime != nil {
					fmt.Printf("Start Time: %s\n", task.Status.StartTime.Format("2006-01-02 15:04:05"))
				}
				if task.Status.CompletionTime != nil {
					fmt.Printf("Completion Time: %s\n", task.Status.CompletionTime.Format("2006-01-02 15:04:05"))
				}
				if task.Status.Message != "" {
					fmt.Printf("Message: %s\n", task.Status.Message)
				}
			} else {
				// Default output format
				fmt.Printf("Name:        %s\n", task.Name)
				fmt.Printf("Namespace:   %s\n", task.Namespace)
				fmt.Printf("Type:        %s\n", task.Spec.Type)
				fmt.Printf("Phase:       %s\n", task.Status.Phase)
				if task.Spec.Model != "" {
					fmt.Printf("Model:       %s\n", task.Spec.Model)
				}
				fmt.Printf("Credentials: %s (secret: %s)\n", task.Spec.Credentials.Type, task.Spec.Credentials.SecretRef.Name)
				fmt.Printf("\nPrompt:\n")
				// Indent the prompt for better readability
				for _, line := range strings.Split(task.Spec.Prompt, "\n") {
					fmt.Printf("  %s\n", line)
				}
				fmt.Println()
				if task.Status.JobName != "" {
					fmt.Printf("Job:         %s\n", task.Status.JobName)
				}
				if task.Status.PodName != "" {
					fmt.Printf("Pod:         %s\n", task.Status.PodName)
				}
				if task.Status.StartTime != nil {
					fmt.Printf("Start Time:  %s\n", task.Status.StartTime.Format("2006-01-02 15:04:05"))
				}
				if task.Status.CompletionTime != nil {
					fmt.Printf("Completion:  %s\n", task.Status.CompletionTime.Format("2006-01-02 15:04:05"))
				}
				if task.Status.Message != "" {
					fmt.Printf("Message:     %s\n", task.Status.Message)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (yaml or json)")

	return cmd
}
