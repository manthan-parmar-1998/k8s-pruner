package resources

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindCompletedPods finds completed pods that are no longer needed
func (d *ResourceDetector) FindCompletedPods(namespace string, olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "Pods",
		Items:        []ResourceItem{},
	}

	// Get all Pods
	pods, err := d.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return result, err
	}

	// Find completed pods
	for _, pod := range pods.Items {
		// Skip pods that are not in Succeeded or Failed phase
		if pod.Status.Phase != corev1.PodSucceeded && pod.Status.Phase != corev1.PodFailed {
			continue
		}

		// Skip pods owned by Jobs or other controllers
		if len(pod.OwnerReferences) > 0 {
			// These will be cleaned up by their controllers
			continue
		}

		// Check age if filter is provided
		if olderThan != nil && pod.CreationTimestamp.Time.After(*olderThan) {
			continue
		}

		result.Items = append(result.Items, ResourceItem{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Age:       pod.CreationTimestamp.Time,
		})
	}

	return result, nil
}
