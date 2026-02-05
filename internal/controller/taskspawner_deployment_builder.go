package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

const (
	// DefaultSpawnerImage is the default image for the spawner binary.
	DefaultSpawnerImage = "gjkim42/axon-spawner:latest"

	// SpawnerServiceAccount is the service account used by spawner Deployments.
	SpawnerServiceAccount = "axon-spawner"

	// SpawnerClusterRole is the ClusterRole referenced by spawner RoleBindings.
	SpawnerClusterRole = "axon-spawner-role"
)

// DeploymentBuilder constructs Kubernetes Deployments for TaskSpawners.
type DeploymentBuilder struct {
	SpawnerImage           string
	SpawnerImagePullPolicy corev1.PullPolicy
}

// NewDeploymentBuilder creates a new DeploymentBuilder.
func NewDeploymentBuilder() *DeploymentBuilder {
	return &DeploymentBuilder{SpawnerImage: DefaultSpawnerImage}
}

// Build creates a Deployment for the given TaskSpawner.
func (b *DeploymentBuilder) Build(ts *axonv1alpha1.TaskSpawner) *appsv1.Deployment {
	replicas := int32(1)

	args := []string{
		"--taskspawner-name=" + ts.Name,
		"--taskspawner-namespace=" + ts.Namespace,
	}

	var envVars []corev1.EnvVar
	if ts.Spec.When.GitHubIssues != nil && ts.Spec.When.GitHubIssues.TokenSecretRef != nil {
		envVars = append(envVars, corev1.EnvVar{
			Name: "GITHUB_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: ts.Spec.When.GitHubIssues.TokenSecretRef.Name,
					},
					Key: "GITHUB_TOKEN",
				},
			},
		})
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "axon",
		"app.kubernetes.io/component":  "spawner",
		"app.kubernetes.io/managed-by": "axon-controller",
		"axon.io/taskspawner":          ts.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ts.Name,
			Namespace: ts.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: SpawnerServiceAccount,
					RestartPolicy:      corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							Name:            "spawner",
							Image:           b.SpawnerImage,
							ImagePullPolicy: b.SpawnerImagePullPolicy,
							Args:            args,
							Env:             envVars,
						},
					},
				},
			},
		},
	}
}
