package argocd

import (
	"fmt"

	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toProject(name string, group string) *appv1.AppProject {
	typeMeta := metav1.TypeMeta{APIVersion: "argoproj.io/v1alpha1", Kind: "AppProject"}
	objectMeta := metav1.ObjectMeta{Name: name, Namespace: env.CloudKnitSystemNamespace()}
	spec := appv1.AppProjectSpec{
		SourceRepos:              []string{"*"},
		Destinations:             []appv1.ApplicationDestination{{Server: "*", Namespace: "*"}},
		ClusterResourceWhitelist: []metav1.GroupKind{{Group: "*", Kind: "*"}},
		Roles: []appv1.ProjectRole{
			{
				Name:   "frontend",
				Groups: []string{fmt.Sprintf("%s:%s", group, name)},
				Policies: []string{
					fmt.Sprintf("p, proj:%s:frontend, applications, get, %s/*, allow", name, name),
					fmt.Sprintf("p, proj:%s:frontend, applications, delete, %s/*, allow", name, name),
					fmt.Sprintf("p, proj:%s:frontend, applications, sync, %s/*, allow", name, name),
				},
			},
		},
	}
	return &appv1.AppProject{TypeMeta: typeMeta, ObjectMeta: objectMeta, Spec: spec, Status: appv1.AppProjectStatus{}}
}

func toBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
