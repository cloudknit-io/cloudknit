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
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/log"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// +k8s:deepcopy-gen=false
type TeamValidator interface {
	ValidateTeamCreate(context.Context, *Team) error
	ValidateTeamUpdate(context.Context, *Team) error
}

var (
	tv TeamValidator
	tl = log.NewLogger().WithFields(logrus.Fields{"name": "controller.TeamValidator"})
)

func (in *Team) SetupWebhookWithManager(mgr ctrl.Manager, vldtr TeamValidator) error {
	tv = vldtr
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

//+kubebuilder:webhook:path=/validate-stable-compuzest-com-v1-team,mutating=false,failurePolicy=fail,sideEffects=None,groups=stable.compuzest.com,resources=teans,verbs=create;update,versions=v1,name=vTeam.kb.io,admissionReviewVersions=v1beta1

var _ webhook.Validator = &Team{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (in *Team) ValidateCreate() error {
	tl.Infof("validating create event for team %s", in.Name)

	return tv.ValidateTeamCreate(ctx, in)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (in *Team) ValidateUpdate(old runtime.Object) error {
	tl.Infof("validate update event for team %s", in.Name)

	return tv.ValidateTeamUpdate(ctx, in)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (in *Team) ValidateDelete() error {
	tl.Infof("validate delete event for team %s", in.Name)

	return nil
}
