package validator

import (
	"context"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TeamValidatorImpl struct {
	kc kClient.Client
	gc git.API
	es eventservice.API
	l  *logrus.Entry
}

func NewTeamValidatorImpl(kc kClient.Client, gc git.API, es eventservice.API, l *logrus.Entry) *TeamValidatorImpl {
	return &TeamValidatorImpl{
		kc: kc,
		es: es,
		gc: gc,
		l:  l,
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

var _ v1.TeamValidator = (*TeamValidatorImpl)(nil)
