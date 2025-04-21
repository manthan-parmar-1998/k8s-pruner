package cmd

import (
	"github.com/spf13/cobra"
)

var (
	namespace  string
	dryRun     bool
	age        string
	context    string
	kubeconfig string
	output     string
	force      bool
	types      []string
	labels     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-pruner",
	Short: "A tool to list and prune unused Kubernetes resources",
	Long: `k8s-pruner helps clean up resources in a Kubernetes cluster to save costs and improve performance.
It can identify and remove unused ConfigMaps, Secrets, PVCs, completed Jobs, and more.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Namespace to target (default is all namespaces)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print resources that would be pruned without actually deleting them")
	rootCmd.PersistentFlags().StringVar(&age, "age", "", "Only consider resources older than this value (e.g., 24h, 7d)")
	rootCmd.PersistentFlags().StringVar(&context, "context", "", "The name of the kubeconfig context to use")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file to use")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Output format (text, json, yaml)")

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(pruneCmd)
	rootCmd.AddCommand(versionCmd)
}
