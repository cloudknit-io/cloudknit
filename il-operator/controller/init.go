package controller

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awscfg"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awseks"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awsssm"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/watcherservices"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/github"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/secrets"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/zlstate"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
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
	secretsClient := awsssm.LazyLoadSSM(awscfg.NewK8sSecretCredentialsLoader(r.Client, env.Config.SharedAWSCredsSecret))

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
	fs := file.NewOSFileService()

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
