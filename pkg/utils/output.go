package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/manthan-parmar-1998/k8s-pruner/pkg/resources"
	"gopkg.in/yaml.v2"
)

// OutputResults outputs the results in the specified format
func OutputResults(results []resources.ResourceList, format string, header string) error {
	if len(results) == 0 {
		fmt.Println("No unused resources found.")
		return nil
	}

	switch strings.ToLower(format) {
	case "json":
		return outputJSON(results)
	case "yaml":
		return outputYAML(results)
	default:
		return outputText(results, header)
	}
}

// outputText outputs results in human-readable text format
func outputText(results []resources.ResourceList, header string) error {
	fmt.Println(header)

	totalCount := 0
	for _, resourceList := range results {
		if len(resourceList.Items) == 0 {
			continue
		}

		fmt.Printf("\n%s:\n", resourceList.ResourceType)
		fmt.Println(strings.Repeat("-", len(resourceList.ResourceType)+1))

		for _, item := range resourceList.Items {
			age := formatAge(time.Since(item.Age))
			fmt.Printf("  %s/%s (age: %s)\n", item.Namespace, item.Name, age)
			totalCount++
		}
	}

	fmt.Printf("\nTotal: %d resources\n", totalCount)
	return nil
}

// outputJSON outputs results in JSON format
func outputJSON(results []resources.ResourceList) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// outputYAML outputs results in YAML format
func outputYAML(results []resources.ResourceList) error {
	yamlData, err := yaml.Marshal(results)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %v", err)
	}

	fmt.Println(string(yamlData))
	return nil
}

// formatAge formats a duration into a human-readable string
func formatAge(d time.Duration) string {
	d = d.Round(time.Minute)

	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour

	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute

	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
