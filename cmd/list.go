package cmd

import (
	"fmt"

	"github.com/manthan-parmar-1998/k8s-pruner/pkg/client"
	"github.com/manthan-parmar-1998/k8s-pruner/pkg/resources"
	"github.com/manthan-parmar-1998/k8s-pruner/pkg/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List unused Kubernetes resources",
	Long: `List unused resources in a Kubernetes cluster.
This command identifies resources that are not being used and can be safely removed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize Kubernetes client
		k8sClient, err := client.NewClient(kubeconfig, context)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %v", err)
		}

		// Parse age filter if provided
		ageFilter, err := utils.ParseAge(age)
		if err != nil {
			return err
		}

		// Create resource detector
		detector := resources.NewResourceDetector(k8sClient)

		// Find unused resources
		results, err := detector.FindAllUnusedResources(namespace, ageFilter, types, labels)
		if err != nil {
			return fmt.Errorf("error finding unused resources: %v", err)
		}

		// Output results
		return utils.OutputResults(results, output, "Found the following unused resources:")
	},
}

func init() {
	listCmd.Flags().StringSliceVar(&types, "types", []string{"configmaps", "secrets", "pvcs", "pods", "jobs", "namespaces"},
		"Resource types to check (configmaps, secrets, pvcs, pods, jobs, namespaces)")
	listCmd.Flags().StringVar(&labels, "labels", "", "Label selector to filter resources")
}
