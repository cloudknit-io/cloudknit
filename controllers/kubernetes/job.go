package controllers

import (
	stablev1alpha1 "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GeneratePreSyncJob(environment stablev1alpha1.Environment) *batchv1.Job {
	jobNamePrefix := environment.Spec.TeamName + "-" + environment.Spec.EnvName

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobNamePrefix + "-presync",
			Namespace: "argocd",
			Annotations: map[string]string{
				"argocd.argoproj.io/hook":               "PreSync",
				"argocd.argoproj.io/hook-delete-policy": "BeforeHookCreation",
			},
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    jobNamePrefix + "-container",
							Image:   "argoproj/argocli:latest",
							Command: []string{"argo", "watch", jobNamePrefix},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}

	return job
}
