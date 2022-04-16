package mocks

import (
	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetMockEnv1(deleted bool) stablev1.Environment {
	deletionTimestamp := metav1.Time{}
	if deleted {
		deletionTimestamp = metav1.Now()
	}
	return stablev1.Environment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Environment",
			APIVersion: "stable.compuzest.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "development",
			Namespace:         "argocd",
			DeletionTimestamp: &deletionTimestamp,
		},
		Spec: stablev1.EnvironmentSpec{
			TeamName:    "design",
			EnvName:     "dev",
			Description: "test",
			Components: []*stablev1.EnvironmentComponent{
				{
					Name: "networking",
					Type: "terraform",
					Module: &stablev1.Module{
						Source: "aws",
						Name:   "vpc",
					},
					VariablesFile: &stablev1.VariablesFile{
						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
						Path:   "dev/networking.tfvars",
					},
				},
				{
					Name:      "rebrand",
					Type:      "terraform",
					DependsOn: []string{"networking"},
					Module: &stablev1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*stablev1.Variable{
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
					Module: &stablev1.Module{
						Source: "aws",
						Name:   "s3-bucket",
					},
					Variables: []*stablev1.Variable{
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

//
// func GetMockEnv2(deleted bool) stablev1.Environment {
//	deletionTimestamp := metav1.Time{}
//	if deleted {
//		deletionTimestamp = metav1.Now()
//	}
//	return stablev1.Environment{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "Environment",
//			APIVersion: "stable.compuzest.com/v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:              "development",
//			Namespace:         "argocd",
//			DeletionTimestamp: &deletionTimestamp,
//		},
//		Spec: stablev1.EnvironmentSpec{
//			TeamName:    "design",
//			EnvName:     "dev",
//			Description: "test",
//			Components: []*stablev1.EnvironmentComponent{
//				{
//					Name: "networking",
//					Type: "terraform",
//					Module: &stablev1.Module{
//						Source: "aws",
//						Name:   "vpc",
//					},
//					VariablesFile: &stablev1.VariablesFile{
//						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
//						Path:   "dev/networking.tfvars",
//					},
//				},
//				{
//					Name:      "rebrand",
//					Type:      "terraform",
//					DependsOn: []string{"networking"},
//					Module: &stablev1.Module{
//						Source: "aws",
//						Name:   "s3-bucket",
//					},
//					Variables: []*stablev1.Variable{
//						{
//							Name:  "bucket",
//							Value: "dev-banners-sandbox",
//						},
//					},
//				},
//			},
//		},
//	}
//}

// func GetMockEnv3(deleted bool) stablev1.Environment {
//	deletionTimestamp := metav1.Time{}
//	if deleted {
//		deletionTimestamp = metav1.Now()
//	}
//	return stablev1.Environment{
//		TypeMeta: metav1.TypeMeta{
//			Kind:       "Environment",
//			APIVersion: "stable.compuzest.com/v1",
//		},
//		ObjectMeta: metav1.ObjectMeta{
//			Name:              "development",
//			Namespace:         "argocd",
//			DeletionTimestamp: &deletionTimestamp,
//		},
//		Spec: stablev1.EnvironmentSpec{
//			TeamName:    "design",
//			EnvName:     "prod",
//			Description: "test",
//			Components: []*stablev1.EnvironmentComponent{
//				{
//					Name: "networking",
//					Type: "terraform",
//					Module: &stablev1.Module{
//						Source: "aws",
//						Name:   "vpc",
//					},
//					VariablesFile: &stablev1.VariablesFile{
//						Source: "git@github.com:zmart-tech-sandbox/zmart-design-team-config.git",
//						Path:   "dev/networking.tfvars",
//					},
//				},
//				{
//					Name:      "rebrand",
//					Type:      "terraform",
//					DependsOn: []string{"networking"},
//					Module: &stablev1.Module{
//						Source: "aws",
//						Name:   "s3-bucket",
//					},
//					Variables: []*stablev1.Variable{
//						{
//							Name:  "bucket",
//							Value: "dev-banners-sandbox",
//						},
//					},
//				},
//			},
//		},
//	}
//}
