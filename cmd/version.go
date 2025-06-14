package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version = "0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of k8s-pruner",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("k8s-pruner version %s\n", Version)
	},
}
