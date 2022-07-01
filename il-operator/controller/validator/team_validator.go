package validator

import (
	"context"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TeamValidatorImpl struct {
	kc kClient.Client
	fs file.API
	gc git.API
	es eventservice.API
	l  *logrus.Entry
}

func NewTeamValidatorImpl(kc kClient.Client, fs file.API, es eventservice.API) *TeamValidatorImpl {
	return &TeamValidatorImpl{
		kc: kc,
		fs: fs,
		es: es,
		l:  logrus.New().WithField("logger", "controller.TeamValidator"),
	}
}

func (v *TeamValidatorImpl) init(ctx context.Context) error {
	watcherServices, err := watcherservices.NewGitHubServices(ctx, v.kc, env.Config.GitHubCompanyOrganization, v.l)
	if err != nil {
		return errors.Wrap(err, "error instantiating watcher services")
	}

	factory := gitfactory.NewFactory(v.kc, v.l)
	var gitOpts gitfactory.Options
	if env.Config.GitHubCompanyAuthMethod == util.AuthModeSSH {
		gitOpts.SSHOptions = &gitfactory.SSHOptions{SecretName: env.Config.GitSSHSecretName, SecretNamespace: env.SystemNamespace()}
	} else {
		gitOpts.GitHubOptions = &gitfactory.GitHubAppOptions{
			GitHubClient:       watcherServices.CompanyGitClient,
			GitHubOrganization: env.Config.GitHubCompanyOrganization,
		}
	}
	gitClient, err := factory.NewGitClient(ctx, &gitOpts)
	if err != nil {
		return errors.Wrap(err, "error instantiating git client")
	}

	v.gc = gitClient

	return nil
}

func (v *TeamValidatorImpl) ValidateTeamCreate(ctx context.Context, t *v1.Team) error {
	if err := v.init(ctx); err != nil {
		v.l.Errorf(errInitTeamValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitTeamValidator))
	}

	var allErrs field.ErrorList

	if verrs := v.validateNames(t); len(verrs) > 0 {
		allErrs = append(allErrs, verrs...)
	}
	if verrs := v.validateConfigRepo(t); verrs != nil {
		allErrs = append(allErrs, verrs...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	for _, e := range allErrs {
		v.l.Warnf("validating webhook error for create team event: %v", e)
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: v1.CRDGroup,
			Kind:  v1.CRDEnvironment,
		},
		t.Name,
		allErrs,
	)
}

func (v *TeamValidatorImpl) ValidateTeamUpdate(ctx context.Context, t *v1.Team) error {
	var allErrs field.ErrorList

	if verrs := v.validateNames(t); len(verrs) > 0 {
		allErrs = append(allErrs, verrs...)
	}
	if verrs := v.validateConfigRepo(t); verrs != nil {
		allErrs = append(allErrs, verrs...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	for _, e := range allErrs {
		v.l.Warnf("validating webhook error for update team event: %v", e)
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: v1.CRDGroup,
			Kind:  v1.CRDEnvironment,
		},
		t.Name,
		allErrs,
	)
}

func (v *TeamValidatorImpl) validateConfigRepo(t *v1.Team) field.ErrorList {
	var allErrs field.ErrorList

	fld := field.NewPath("spec").Child("overlayFiles")
	if t.Spec.ConfigRepo == nil {
		fe := field.NotFound(fld, "team config repo block must be defined")
		allErrs = append(allErrs, fe)
		return allErrs
	}

	return checkPaths(v.fs, v.gc, t.Spec.ConfigRepo.Source, []string{t.Spec.ConfigRepo.Path}, fld, v.l)
}

func (v *TeamValidatorImpl) validateNames(t *v1.Team) field.ErrorList {
	var allErrs field.ErrorList

	if err := validateRFC1035String(t.Spec.TeamName); err != nil {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, t.Spec.TeamName, err.Error()))
	}
	if err := validateStringLength(t.Spec.TeamName); err != nil {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, t.Spec.TeamName, err.Error()))
	}

	return allErrs
}

var _ v1.TeamValidator = (*TeamValidatorImpl)(nil)
