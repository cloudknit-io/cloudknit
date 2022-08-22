/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package argocd

import (
	"fmt"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toBlankProject(project string, namespace string) *appv1.AppProject {
	typeMeta := metav1.TypeMeta{APIVersion: "argoproj.io/v1alpha1", Kind: "AppProject"}
	objectMeta := metav1.ObjectMeta{Name: project, Namespace: namespace}
	spec := appv1.AppProjectSpec{
		SourceRepos:              []string{"*"},
		Destinations:             []appv1.ApplicationDestination{{Server: "*", Namespace: "*"}},
		ClusterResourceWhitelist: []metav1.GroupKind{{Group: "*", Kind: "*"}},
	}
	return &appv1.AppProject{TypeMeta: typeMeta, ObjectMeta: objectMeta, Spec: spec, Status: appv1.AppProjectStatus{}}
}

func toProject(project, namespace, group string) *appv1.AppProject {
	typeMeta := metav1.TypeMeta{APIVersion: "argoproj.io/v1alpha1", Kind: "AppProject"}
	objectMeta := metav1.ObjectMeta{Name: project, Namespace: namespace}
	spec := appv1.AppProjectSpec{
		SourceRepos:              []string{"*"},
		Destinations:             []appv1.ApplicationDestination{{Server: "*", Namespace: "*"}},
		ClusterResourceWhitelist: []metav1.GroupKind{{Group: "*", Kind: "*"}},
		Roles: []appv1.ProjectRole{
			{
				Name:   "frontend",
				Groups: []string{fmt.Sprintf("%s:%s", group, project)},
				Policies: []string{
					fmt.Sprintf("p, proj:%s:frontend, applications, get, %s/*, allow", project, project),
					fmt.Sprintf("p, proj:%s:frontend, applications, delete, %s/*, allow", project, project),
					fmt.Sprintf("p, proj:%s:frontend, applications, sync, %s/*, allow", project, project),
				},
			},
		},
	}
	return &appv1.AppProject{TypeMeta: typeMeta, ObjectMeta: objectMeta, Spec: spec, Status: appv1.AppProjectStatus{}}
}

func toBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
