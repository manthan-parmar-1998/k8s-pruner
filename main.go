// main.go
package main

import (
	"fmt"
	"os"

	"github.com/manthan-parmar-1998/k8s-pruner/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
