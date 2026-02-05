package controller

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	axonv1alpha1 "github.com/gjkim/axon/api/v1alpha1"
)

const (
	// ClaudeCodeImage is the default image for Claude Code agent.
	ClaudeCodeImage = "gjkim42/claude-code:latest"

	// AgentTypeClaudeCode is the agent type for Claude Code.
	AgentTypeClaudeCode = "claude-code"
)

// JobBuilder constructs Kubernetes Jobs for Tasks.
type JobBuilder struct{}

// NewJobBuilder creates a new JobBuilder.
func NewJobBuilder() *JobBuilder {
	return &JobBuilder{}
}

// Build creates a Job for the given Task.
func (b *JobBuilder) Build(task *axonv1alpha1.Task) (*batchv1.Job, error) {
	switch task.Spec.Type {
	case AgentTypeClaudeCode:
		return b.buildClaudeCodeJob(task)
	default:
		return nil, fmt.Errorf("unsupported agent type: %s", task.Spec.Type)
	}
}

// buildClaudeCodeJob creates a Job for Claude Code agent.
func (b *JobBuilder) buildClaudeCodeJob(task *axonv1alpha1.Task) (*batchv1.Job, error) {
	args := []string{
		"--dangerously-skip-permissions",
		"-p", task.Spec.Prompt,
	}

	if task.Spec.Model != "" {
		args = append(args, "--model", task.Spec.Model)
	}

	var envVars []corev1.EnvVar

	switch task.Spec.Credentials.Type {
	case axonv1alpha1.CredentialTypeAPIKey:
		envVars = append(envVars, corev1.EnvVar{
			Name: "ANTHROPIC_API_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: task.Spec.Credentials.SecretRef.Name,
					},
					Key: "ANTHROPIC_API_KEY",
				},
			},
		})
	case axonv1alpha1.CredentialTypeOAuth:
		envVars = append(envVars, corev1.EnvVar{
			Name: "CLAUDE_CODE_OAUTH_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: task.Spec.Credentials.SecretRef.Name,
					},
					Key: "CLAUDE_CODE_OAUTH_TOKEN",
				},
			},
		})
	}

	backoffLimit := int32(0)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      task.Name,
			Namespace: task.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "axon",
				"app.kubernetes.io/component":  "task",
				"app.kubernetes.io/managed-by": "axon-controller",
				"axon.io/task":                 task.Name,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":       "axon",
						"app.kubernetes.io/component":  "task",
						"app.kubernetes.io/managed-by": "axon-controller",
						"axon.io/task":                 task.Name,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "claude-code",
							Image: ClaudeCodeImage,
							Args:  args,
							Env:   envVars,
						},
					},
				},
			},
		},
	}

	return job, nil
}
