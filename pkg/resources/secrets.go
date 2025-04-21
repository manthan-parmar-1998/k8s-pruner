package resources

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindUnusedSecrets finds Secrets that are not mounted in any Pod
func (d *ResourceDetector) FindUnusedSecrets(namespace string, olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "Secrets",
		Items:        []ResourceItem{},
	}

	// Get all Secrets
	secrets, err := d.client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
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

	// Track which Secrets are in use
	usedSecrets := make(map[string]bool)

	// Check Secret usage in Pods
	for _, pod := range pods.Items {
		// Check volume mounts
		for _, volume := range pod.Spec.Volumes {
			if volume.Secret != nil {
				key := pod.Namespace + "/" + volume.Secret.SecretName
				usedSecrets[key] = true
			}
		}

		// Check environment variables
		for _, container := range pod.Spec.Containers {
			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
					key := pod.Namespace + "/" + env.ValueFrom.SecretKeyRef.Name
					usedSecrets[key] = true
				}
			}

			// Check envFrom
			for _, envFrom := range container.EnvFrom {
				if envFrom.SecretRef != nil {
					key := pod.Namespace + "/" + envFrom.SecretRef.Name
					usedSecrets[key] = true
				}
			}
		}

		// Check image pull secrets
		for _, pullSecret := range pod.Spec.ImagePullSecrets {
			key := pod.Namespace + "/" + pullSecret.Name
			usedSecrets[key] = true
		}
	}

	// Get all ServiceAccounts to check for token secrets
	serviceAccounts, err := d.client.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	// Check Secret usage in ServiceAccounts
	for _, sa := range serviceAccounts.Items {
		for _, secret := range sa.Secrets {
			key := sa.Namespace + "/" + secret.Name
			usedSecrets[key] = true
		}
		if sa.ImagePullSecrets != nil {
			for _, pullSecret := range sa.ImagePullSecrets {
				key := sa.Namespace + "/" + pullSecret.Name
				usedSecrets[key] = true
			}
		}
	}

	// Add unused Secrets to result
	for _, secret := range secrets.Items {
		key := secret.Namespace + "/" + secret.Name

		// Skip service account tokens, TLS secrets, and system secrets
		if secret.Type == "kubernetes.io/service-account-token" ||
			secret.Type == "kubernetes.io/tls" ||
			secret.Namespace == "kube-system" {
			continue
		}

		// Check if Secret is unused
		if !usedSecrets[key] {
			// Check age if filter is provided
			if olderThan != nil && secret.CreationTimestamp.Time.After(*olderThan) {
				continue
			}

			result.Items = append(result.Items, ResourceItem{
				Name:      secret.Name,
				Namespace: secret.Namespace,
				Age:       secret.CreationTimestamp.Time,
			})
		}
	}

	return result, nil
}
