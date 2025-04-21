package resources

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindUnusedPVCs finds PVCs that are not mounted in any Pod
func (d *ResourceDetector) FindUnusedPVCs(namespace string, olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "PersistentVolumeClaims",
		Items:        []ResourceItem{},
	}

	// Get all PVCs
	pvcs, err := d.client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return result, err
	}

	// Get all Pods
	pods, err := d.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	// Track which PVCs are in use
	usedPVCs := make(map[string]bool)

	// Check PVC usage in Pods
	for _, pod := range pods.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				key := pod.Namespace + "/" + volume.PersistentVolumeClaim.ClaimName
				usedPVCs[key] = true
			}
		}
	}

	// Add unused PVCs to result
	for _, pvc := range pvcs.Items {
		key := pvc.Namespace + "/" + pvc.Name

		// Check if PVC is unused
		if !usedPVCs[key] {
			// Check age if filter is provided
			if olderThan != nil && pvc.CreationTimestamp.Time.After(*olderThan) {
				continue
			}

			result.Items = append(result.Items, ResourceItem{
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				Age:       pvc.CreationTimestamp.Time,
			})
		}
	}

	return result, nil
}
