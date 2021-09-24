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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var environmentlog = logf.Log.WithName("environment-resource")

func (in *Environment) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

//+kubebuilder:webhook:path=/validate-stable-compuzest-com-v1-environment,mutating=false,failurePolicy=fail,sideEffects=None,groups=stable.compuzest.com,resources=environments,verbs=create;update,versions=v1,name=venvironment.kb.io,admissionReviewVersions=v1beta1

var _ webhook.Validator = &Environment{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *Environment) ValidateCreate() error {
	environmentlog.Info("validating environment create", "name", in.Name)

	return ValidateEnvironmentCreate(in)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *Environment) ValidateUpdate(old runtime.Object) error {
	environmentlog.Info("validate environment update", "name", in.Name)

	return ValidateEnvironmentUpdate(in)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *Environment) ValidateDelete() error {
	environmentlog.Info("validate environment delete", "name", in.Name)

	return nil
}
