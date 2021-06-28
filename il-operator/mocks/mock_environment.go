package mocks

import (
	"github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetMockEnv1(deleted bool) v1alpha1.Environment {
	deletionTimestamp := metav1.Time{}
	if deleted {
		deletionTimestamp = metav1.Now()
	}
	return v1alpha1.Environment{
		TypeMeta: v1.TypeMeta{
			Kind:       "Environment",
			APIVersion: "stable.compuzest.com/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:              "development",
			Namespace:         "argocd",
			DeletionTimestamp: &deletionTimestamp,
		},
		Spec: v1alpha1.EnvironmentSpec{
			TeamName:    "design",
			EnvName:     "dev",
			Description: "test",
			EnvironmentComponent: []*v1alpha1.EnvironmentComponent{
				{
					Name: "networking",
					Type: "terraform",
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "vpc",
					},
					VariablesFile: &v1alpha1.VariablesFile{
						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
						Path:   "dev/networking.tfvars",
					},
				},
				{
					Name:      "rebrand",
					Type:      "terraform",
					DependsOn: []string{"networking"},
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*v1alpha1.Variable{
						{
							Name:  "bucket",
							Value: "dev-banners-sandbox",
						},
					},
				},
				{
					Name:      "overlay",
					Type:      "terraform",
					DependsOn: []string{"networking", "rebrand"},
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*v1alpha1.Variable{
						{
							Name:  "bucket",
							Value: "dev-overlays-sandbox",
						},
					},
				},
			},
		},
	}
}

func GetMockEnv2(deleted bool) v1alpha1.Environment {
	deletionTimestamp := metav1.Time{}
	if deleted {
		deletionTimestamp = metav1.Now()
	}
	return v1alpha1.Environment{
		TypeMeta: v1.TypeMeta{
			Kind:       "Environment",
			APIVersion: "stable.compuzest.com/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:              "development",
			Namespace:         "argocd",
			DeletionTimestamp: &deletionTimestamp,
		},
		Spec: v1alpha1.EnvironmentSpec{
			TeamName:    "design",
			EnvName:     "dev",
			Description: "test",
			EnvironmentComponent: []*v1alpha1.EnvironmentComponent{
				{
					Name: "networking",
					Type: "terraform",
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "vpc",
					},
					VariablesFile: &v1alpha1.VariablesFile{
						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
						Path:   "dev/networking.tfvars",
					},
				},
				{
					Name:      "rebrand",
					Type:      "terraform",
					DependsOn: []string{"networking"},
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*v1alpha1.Variable{
						{
							Name:  "bucket",
							Value: "dev-banners-sandbox",
						},
					},
				},
			},
		},
	}
}

func GetMockEnv3(deleted bool) v1alpha1.Environment {
	deletionTimestamp := metav1.Time{}
	if deleted {
		deletionTimestamp = metav1.Now()
	}
	return v1alpha1.Environment{
		TypeMeta: v1.TypeMeta{
			Kind:       "Environment",
			APIVersion: "stable.compuzest.com/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:              "development",
			Namespace:         "argocd",
			DeletionTimestamp: &deletionTimestamp,
		},
		Spec: v1alpha1.EnvironmentSpec{
			TeamName:    "design",
			EnvName:     "prod",
			Description: "test",
			EnvironmentComponent: []*v1alpha1.EnvironmentComponent{
				{
					Name: "networking",
					Type: "terraform",
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "vpc",
					},
					VariablesFile: &v1alpha1.VariablesFile{
						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
						Path:   "dev/networking.tfvars",
					},
				},
				{
					Name:      "rebrand",
					Type:      "terraform",
					DependsOn: []string{"networking"},
					Module: &v1alpha1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*v1alpha1.Variable{
						{
							Name:  "bucket",
							Value: "dev-banners-sandbox",
						},
					},
				},
			},
		},
	}
}
