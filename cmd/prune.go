package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manthan-parmar-1998/k8s-pruner/pkg/client"
	"github.com/manthan-parmar-1998/k8s-pruner/pkg/resources"
	"github.com/manthan-parmar-1998/k8s-pruner/pkg/utils"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune unused Kubernetes resources",
	Long: `Prune (delete) unused resources in a Kubernetes cluster.
This command removes resources that are not being used to free up cluster resources.`,
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

		// Count total resources
		totalCount := 0
		for _, resourceList := range results {
			totalCount += len(resourceList.Items)
		}

		if totalCount == 0 {
			fmt.Println("No unused resources found.")
			return nil
		}

		// Output what we found
		if err := utils.OutputResults(results, output, "Found the following unused resources:"); err != nil {
			return err
		}

		// If dry run, exit here
		if dryRun {
			fmt.Println("\nDRY RUN: No resources were pruned.")
			return nil
		}

		// Confirm deletion unless force flag is set
		if !force {
			fmt.Printf("\nAre you sure you want to delete these %d resources? (y/N): ", totalCount)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %v", err)
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Operation cancelled.")
				return nil
			}
		}

		// Delete resources
		deleted, err := detector.DeleteUnusedResources(results)
		if err != nil {
			return fmt.Errorf("error deleting resources: %v", err)
		}

		fmt.Printf("Successfully pruned %d resources.\n", deleted)
		return nil
	},
}

func init() {
	pruneCmd.Flags().StringSliceVar(&types, "types", []string{"configmaps", "secrets", "pvcs", "pods", "jobs", "namespaces"},
		"Resource types to prune (configmaps, secrets, pvcs, pods, jobs, namespaces)")
	pruneCmd.Flags().StringVar(&labels, "labels", "", "Label selector to filter resources")
	pruneCmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt before deleting resources")
}
