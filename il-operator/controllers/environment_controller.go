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

	github "github.com/compuzest/environment-operator/controllers/github"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	stablev1alpha1 "github.com/compuzest/environment-operator/api/v1alpha1"
	argocd "github.com/compuzest/environment-operator/controllers/argocd"
	argoWorkflow "github.com/compuzest/environment-operator/controllers/argoworkflow"
	fileutil "github.com/compuzest/environment-operator/controllers/file"
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
	ctx := context.Background()
	//log := r.Log.WithValues("environment", req.NamespacedName)

	environment := &stablev1alpha1.Environment{}

	r.Get(ctx, req.NamespacedName, environment)

	envPrefix := environment.Spec.CustomerId + "/" + environment.Spec.Name

	env := argocd.GenerateEnvironmentApp(*environment)
	fileutil.SaveYamlFile(*env, envPrefix+".yaml")

	for _, terraformConfig := range environment.Spec.TerraformConfigs {
		if terraformConfig.Variables != nil {
			filePath := environment.Spec.CustomerId + "/" + environment.Spec.Name + "/" + terraformConfig.Name + ".tfvars"
			fileutil.SaveVarsToFile(terraformConfig.Variables, filePath)
			terraformConfig.VariablesFile = &stablev1alpha1.VariablesFile{
				Source: "git@github.com:CompuZest/terraform-environment.git",
				Path:   filePath,
			}
		}

		application := argocd.GenerateTerraformConfigApps(*environment, *terraformConfig)

		fileutil.SaveYamlFile(*application, envPrefix+"/"+terraformConfig.Name+".yaml")
	}

	workflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	fileutil.SaveYamlFile(*workflow, envPrefix+"/wofw.yaml")

	github.CommitAndPushFiles("CompuZest", "terraform-environment", environment.Spec.CustomerId+"/", "master", "Adarsh Shah", "shahadarsh@gmail.com")

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1alpha1.Environment{}).
		Complete(r)
}
