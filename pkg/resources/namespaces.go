// pkg/resources/namespaces.go
package resources

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindUnusedNamespaces finds namespaces that don't contain any resources
func (d *ResourceDetector) FindUnusedNamespaces(olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "Namespaces",
		Items:        []ResourceItem{},
	}

	// Get all namespaces
	namespaces, err := d.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return result, err
	}

	// Check each namespace for resources
	for _, ns := range namespaces.Items {
		// Skip system and default namespaces
		if isSystemNamespace(ns.Name) {
			continue
		}

		// Check age if filter is provided
		if olderThan != nil && ns.CreationTimestamp.Time.After(*olderThan) {
			continue
		}

		// Check if namespace has any resources
		isEmpty, err := d.isNamespaceEmpty(ns.Name)
		if err != nil {
			return result, err
		}

		if isEmpty {
			result.Items = append(result.Items, ResourceItem{
				Name:      ns.Name,
				Namespace: "",
				Age:       ns.CreationTimestamp.Time,
			})
		}
	}

	return result, nil
}

// isNamespaceEmpty checks if a namespace has any resources
func (d *ResourceDetector) isNamespaceEmpty(namespace string) (bool, error) {
	ctx := context.Background()

	// Check for pods
	pods, err := d.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(pods.Items) > 0 {
		return false, nil
	}

	// Check for services
	services, err := d.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	// Skip the default kubernetes service
	if len(services.Items) > 1 || (len(services.Items) == 1 && services.Items[0].Name != "kubernetes") {
		return false, nil
	}

	// Check for deployments
	deployments, err := d.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(deployments.Items) > 0 {
		return false, nil
	}

	// Check for statefulsets
	statefulsets, err := d.client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(statefulsets.Items) > 0 {
		return false, nil
	}

	// Check for daemonsets
	daemonsets, err := d.client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(daemonsets.Items) > 0 {
		return false, nil
	}

	// Check for configmaps (excluding default ones)
	configmaps, err := d.client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, cm := range configmaps.Items {
		if !isDefaultConfigMap(cm.Name) {
			return false, nil
		}
	}

	// Check for secrets (excluding default ones)
	secrets, err := d.client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, secret := range secrets.Items {
		if !isDefaultSecret(secret.Name, secret.Type) {
			return false, nil
		}
	}

	// If we get here, the namespace is empty
	return true, nil
}

// isSystemNamespace returns true if the namespace is a system namespace that should not be pruned
func isSystemNamespace(name string) bool {
	systemNamespaces := map[string]bool{
		"default":                     true,
		"kube-system":                 true,
		"kube-public":                 true,
		"kube-node-lease":             true,
		"cert-manager":                true,
		"ingress-nginx":               true,
		"monitoring":                  true,
		"istio-system":                true,
		"knative-serving":             true,
		"cattle-system":               true,
		"cattle-global-data":          true,
		"cattle-global-nt":            true,
		"cattle-impersonation-system": true,
		"local-path-storage":          true,
	}
	return systemNamespaces[name]
}

// isDefaultConfigMap returns true if the ConfigMap is a default one that should be ignored
func isDefaultConfigMap(name string) bool {
	defaultConfigMaps := map[string]bool{
		"kube-root-ca.crt": true,
	}
	return defaultConfigMaps[name]
}

// isDefaultSecret returns true if the Secret is a default one that should be ignored
// secretType is of type corev1.SecretType
func isDefaultSecret(name string, secretType corev1.SecretType) bool {
	if secretType == corev1.SecretTypeServiceAccountToken {
		return true
	}

	defaultSecrets := map[string]bool{
		"default-token": true,
	}

	// Check if the name starts with any of the default prefixes
	for defaultName := range defaultSecrets {
		if len(name) >= len(defaultName) && name[:len(defaultName)] == defaultName {
			return true
		}
	}

	return false
}
