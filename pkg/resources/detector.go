package resources

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ResourceItem represents a Kubernetes resource
type ResourceItem struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Age       time.Time `json:"age"`
}

// ResourceList represents a list of resources of a specific type
type ResourceList struct {
	ResourceType string         `json:"resourceType"`
	Items        []ResourceItem `json:"items"`
}

// ResourceDetector handles detection of unused resources
type ResourceDetector struct {
	client kubernetes.Interface
}

// NewResourceDetector creates a new ResourceDetector
func NewResourceDetector(client kubernetes.Interface) *ResourceDetector {
	return &ResourceDetector{client: client}
}

// FindAllUnusedResources finds all unused resources of the specified types
func (d *ResourceDetector) FindAllUnusedResources(namespace string, olderThan *time.Time, types []string, labelSelector string) ([]ResourceList, error) {
	var results []ResourceList

	for _, resourceType := range types {
		var resourceList ResourceList
		var err error

		switch resourceType {
		case "configmaps":
			resourceList, err = d.FindUnusedConfigMaps(namespace, olderThan, labelSelector)
		case "secrets":
			resourceList, err = d.FindUnusedSecrets(namespace, olderThan, labelSelector)
		case "pvcs":
			resourceList, err = d.FindUnusedPVCs(namespace, olderThan, labelSelector)
		case "pods":
			resourceList, err = d.FindCompletedPods(namespace, olderThan, labelSelector)
		case "jobs":
			resourceList, err = d.FindCompletedJobs(namespace, olderThan, labelSelector)
		case "namespaces":
			if namespace == "" {
				resourceList, err = d.FindUnusedNamespaces(olderThan, labelSelector)
			}
		}

		if err != nil {
			return nil, err
		}

		if len(resourceList.Items) > 0 {
			results = append(results, resourceList)
		}
	}

	return results, nil
}

// DeleteUnusedResources deletes the specified unused resources
func (d *ResourceDetector) DeleteUnusedResources(resources []ResourceList) (int, error) {
	deletedCount := 0
	ctx := context.Background() // Create a context

	for _, resourceList := range resources {
		for _, item := range resourceList.Items {
			var err error

			switch resourceList.ResourceType {
			case "ConfigMaps":
				err = d.client.CoreV1().ConfigMaps(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{})
			case "Secrets":
				err = d.client.CoreV1().Secrets(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{})
			case "PersistentVolumeClaims":
				err = d.client.CoreV1().PersistentVolumeClaims(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{})
			case "Pods":
				err = d.client.CoreV1().Pods(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{})
			case "Jobs":
				err = d.client.BatchV1().Jobs(item.Namespace).Delete(ctx, item.Name, metav1.DeleteOptions{})
			case "Namespaces":
				err = d.client.CoreV1().Namespaces().Delete(ctx, item.Name, metav1.DeleteOptions{})
			}

			if err != nil {
				return deletedCount, err
			}

			deletedCount++
		}
	}

	return deletedCount, nil
}
