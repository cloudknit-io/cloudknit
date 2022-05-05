package controllers

import (
	"context"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/aws/awscfg"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/aws/awseks"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/aws/awsssm"
	"github.com/compuzest/zlifecycle-il-operator/controllers/lib/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controllers/lib/watcherservices"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/il"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controllers/codegen/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/zlstate"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	perrors "github.com/pkg/errors"
)

type EnvironmentServices struct {
	ZLStateManagerClient zlstate.API
	ArgocdClient         argocd.API
	ArgoWorkflowClient   argoworkflow.API
	WatcherServices      *watcherservices.WatcherServices
	SecretsClient        secrets.API
	K8sClient            awseks.API
	ILService            *il.Service
	CompanyGitClient     git.API
	FileService          file.API
}

type Tokens struct {
	GitILToken string
}

func (r *EnvironmentReconciler) initServices(ctx context.Context, environment *v1.Environment) (*EnvironmentServices, error) {
	zlstateManagerClient := zlstate.NewHTTPStateManager(ctx, r.LogV2)
	argocdClient := argocd.NewHTTPClient(ctx, r.LogV2, env.Config.ArgocdServerURL)
	argoworkflowClient := argoworkflow.NewHTTPClient(ctx, env.Config.ArgoWorkflowsServerURL)
	watcherServices, err := watcherservices.NewGitHubServices(ctx, r.Client, env.Config.GitHubCompanyOrganization, r.LogV2)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating watcher services")
	}
	secretsClient := awsssm.LazyLoadSSM(awscfg.NewK8sSecretCredentialsLoader(r.Client))

	secretsMeta := secrets2.Identifier{
		Company:     env.Config.CompanyName,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
	}
	cl := awscfg.NewSSMCredentialsLoader(secretsClient, &secretsMeta, r.LogV2)
	k8sClient := awseks.LazyLoadEKS(ctx, cl, r.LogV2)

	ilToken, err := github.GenerateInstallationToken(r.LogV2, watcherServices.ILGitClient, env.Config.GitILRepositoryOwner)
	if err != nil {
		return nil, perrors.Wrapf(err, "error generating installation token for git organization [%s]", env.Config.GitILRepositoryOwner)
	}
	ilService, err := il.NewService(ctx, ilToken, r.LogV2)
	if err != nil {
		return nil, perrors.Wrap(err, "error getting environment from k8s cache")
	}

	var companyGitClient git.API

	factory := gitfactory.NewFactory(r.Client, r.LogV2)
	var gitOpts gitfactory.Options
	if env.Config.GitHubCompanyAuthMethod == util.AuthModeSSH {
		gitOpts.SSHOptions = &gitfactory.SSHOptions{SecretName: env.Config.GitSSHSecretName, SecretNamespace: env.SystemNamespace()}
	} else {
		gitOpts.GitHubOptions = &gitfactory.GitHubAppOptions{
			GitHubClient:       watcherServices.CompanyGitClient,
			GitHubOrganization: env.Config.GitHubCompanyOrganization,
		}
	}
	companyGitClient, err = factory.NewGitClient(ctx, &gitOpts)
	if err != nil {
		return nil, perrors.Wrap(err, "error instantiating git client")
	}
	fs := file.NewOsFileService()

	return &EnvironmentServices{
		ZLStateManagerClient: zlstateManagerClient,
		ArgocdClient:         argocdClient,
		ArgoWorkflowClient:   argoworkflowClient,
		WatcherServices:      watcherServices,
		SecretsClient:        secretsClient,
		K8sClient:            k8sClient,
		ILService:            ilService,
		CompanyGitClient:     companyGitClient,
		FileService:          fs,
	}, nil
}
