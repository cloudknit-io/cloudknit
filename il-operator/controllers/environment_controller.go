/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"

	"fmt"
	"io/ioutil"
	//"k8s.io/client-go/kubernetes/scheme"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/scheme"
	//workflowv1alpha1 "github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	v1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"

	"k8s.io/client-go/rest"
)

// EnvironmentReconciler reconciles a Environment object
type EnvironmentReconciler struct {
	kClient.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments/status,verbs=get;update;patch

func (r *EnvironmentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("environment", req.NamespacedName)

	var err error
	var content []byte

	content, err = ioutil.ReadFile("dev.yaml")

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, groupVersionKind, err := decode(content, nil, nil)

	fmt.Printf("%#v\n", groupVersionKind)

	if err != nil {
		fmt.Printf("%#v\n", (fmt.Sprintf("Error while decoding YAML object. Err was: %s", err)))
	}

	workflowObj := obj.(*v1alpha1.Workflow)

	fmt.Printf("%#v\n", workflowObj)
	config, err := rest.InClusterConfig()

	wfClient := wfclientset.NewForConfigOrDie(config).ArgoprojV1alpha1().Workflows("argo")

	createdWf, err := wfClient.Create(workflowObj)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", createdWf.Name)
	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}
