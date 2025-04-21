package resources

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FindCompletedJobs finds completed jobs that are no longer needed
func (d *ResourceDetector) FindCompletedJobs(namespace string, olderThan *time.Time, labelSelector string) (ResourceList, error) {
	ctx := context.Background()
	result := ResourceList{
		ResourceType: "Jobs",
		Items:        []ResourceItem{},
	}

	// Get all Jobs
	jobs, err := d.client.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return result, err
	}

	// Find completed jobs
	for _, job := range jobs.Items {
		// Skip jobs that are not completed
		if job.Status.CompletionTime == nil {
			continue
		}

		// Skip jobs owned by CronJobs
		isCronJobChild := false
		for _, owner := range job.OwnerReferences {
			if owner.Kind == "CronJob" {
				isCronJobChild = true
				break
			}
		}
		if isCronJobChild {
			continue
		}

		// Check age if filter is provided
		if olderThan != nil && job.CreationTimestamp.Time.After(*olderThan) {
			continue
		}

		result.Items = append(result.Items, ResourceItem{
			Name:      job.Name,
			Namespace: job.Namespace,
			Age:       job.CreationTimestamp.Time,
		})
	}

	return result, nil
}
