package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

// NewListCmd creates a new list command
func NewListCmd() *cobra.Command {
	var allNamespaces bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Long:  `List all tasks in the namespace.`,
		Example: `  # List tasks in default namespace
  axonctl list

  # List tasks in specific namespace
  axonctl list -n my-namespace

  # List tasks in all namespaces
  axonctl list --all-namespaces`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAxonClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			ctx := context.Background()
			taskList := &axonv1alpha1.TaskList{}

			if !allNamespaces {
				if err := client.List(ctx, taskList, &ctrlclient.ListOptions{Namespace: namespace}); err != nil {
					return fmt.Errorf("failed to list tasks: %w", err)
				}
			} else {
				if err := client.List(ctx, taskList); err != nil {
					return fmt.Errorf("failed to list tasks: %w", err)
				}
			}

			if len(taskList.Items) == 0 {
				fmt.Println("No tasks found")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			if allNamespaces {
				fmt.Fprintln(w, "NAMESPACE\tNAME\tTYPE\tPHASE\tAGE")
			} else {
				fmt.Fprintln(w, "NAME\tTYPE\tPHASE\tAGE")
			}

			for _, task := range taskList.Items {
				age := metav1.Now().Sub(task.CreationTimestamp.Time).Round(time.Second)
				if allNamespaces {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
						task.Namespace, task.Name, task.Spec.Type, task.Status.Phase, age)
				} else {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
						task.Name, task.Spec.Type, task.Status.Phase, age)
				}
			}
			w.Flush()

			return nil
		},
	}

	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List tasks from all namespaces")

	return cmd
}
