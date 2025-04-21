package resources

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindUnusedConfigMaps finds ConfigMaps that are not mounted in any Pod
func (d *ResourceDetector) FindUnusedConfigMaps(namespace string, olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "ConfigMaps",
		Items:        []ResourceItem{},
	}

	// Get all ConfigMaps
	configMaps, err := d.client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{
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

	// Track which ConfigMaps are in use
	usedConfigMaps := make(map[string]bool)

	// Check ConfigMap usage in Pods
	for _, pod := range pods.Items {
		// Check volume mounts
		for _, volume := range pod.Spec.Volumes {
			if volume.ConfigMap != nil {
				key := pod.Namespace + "/" + volume.ConfigMap.Name
				usedConfigMaps[key] = true
			}
		}

		// Check environment variables
		for _, container := range pod.Spec.Containers {
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
					key := pod.Namespace + "/" + env.ValueFrom.ConfigMapKeyRef.Name
					usedConfigMaps[key] = true
				}
			}

			// Check envFrom
			for _, envFrom := range container.EnvFrom {
				if envFrom.ConfigMapRef != nil {
					key := pod.Namespace + "/" + envFrom.ConfigMapRef.Name
					usedConfigMaps[key] = true
				}
			}
		}
	}

	// Add unused ConfigMaps to result
	for _, cm := range configMaps.Items {
		key := cm.Namespace + "/" + cm.Name

		// Skip kube-system ConfigMaps and those with special prefixes
		if cm.Namespace == "kube-system" ||
			cm.Name == "kube-root-ca.crt" {
			continue
		}

		// Check if ConfigMap is unused
		if !usedConfigMaps[key] {
			// Check age if filter is provided
			if olderThan != nil && cm.CreationTimestamp.Time.After(*olderThan) {
				continue
			}

			result.Items = append(result.Items, ResourceItem{
				Name:      cm.Name,
				Namespace: cm.Namespace,
				Age:       cm.CreationTimestamp.Time,
			})
		}
	}

	return result, nil
}
