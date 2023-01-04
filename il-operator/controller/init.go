package controller

import (
	"context"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"

	argoworkflow2 "github.com/compuzest/zlifecycle-il-operator/controller/common/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/aws/awscfg"
	awseks2 "github.com/compuzest/zlifecycle-il-operator/controller/common/aws/awseks"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/aws/awsssm"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/cloudknitservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/statemanager"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/github"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"
	secrets2 "github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	"github.com/pkg/errors"
)

type EnvironmentServices struct {
	CloudKnitServiceClient cloudknitservice.API
	StateManagerClient     statemanager.API
	EventService           eventservice.API
	ArgocdClient           argocd.API
	ArgoWorkflowClient     argoworkflow2.API
	WatcherServices        *watcherservices.WatcherServices
	SecretsClient          secret.API
	K8sClient              awseks2.API
	ILService              *il.Service
	CompanyGitClient       git.API
	FileService            file.API
}

type Tokens struct {
	GitILToken string
}

func (r *EnvironmentReconciler) initServices(ctx context.Context, environment *v1.Environment) (*EnvironmentServices, error) {
	stateManagerClient := statemanager.NewService(env.Config.ZLifecycleStateManagerURL)
	cloudKnitServiceClient := cloudknitservice.NewService(env.Config.ZLifecycleAPIURL)
	eventServiceClient := eventservice.NewService(env.Config.ZLifecycleEventServiceURL)
	argocdClient := argocd.NewHTTPClient(ctx, r.LogV2, env.Config.ArgocdServerURL)
	argoworkflowClient := argoworkflow2.NewHTTPClient(ctx, env.Config.ArgoWorkflowsServerURL)

	organization, err := cloudKnitServiceClient.GetOrganization(ctx, env.Config.CompanyName, r.LogV2)

	if err != nil {
		return nil, errors.Wrap(err, "error getting Organization Response")
	}

	if len(organization.GitHubOrgName) != 0 {
		env.Config.GitHubCompanyOrganization = organization.GitHubOrgName
	}
	if len(organization.GitHubRepo) != 0 {
		env.Config.GitHubRepoURL = organization.GitHubRepo
	}

	watcherServices, err := watcherservices.NewGitHubServices(ctx, r.Client, env.Config.GitHubCompanyOrganization, r.LogV2)
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating watcher services")
	}
	secretsClient := awsssm.LazyLoadSSM(awscfg.NewK8sSecretCredentialsLoader(r.Client, env.Config.SharedAWSCredsSecret))

	secretsMeta := secrets2.Identifier{
		Company:     env.Config.CompanyName,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
	}
	cl := awscfg.NewSSMCredentialsLoader(secretsClient, &secretsMeta, r.LogV2)
	k8sClient := awseks2.LazyLoadEKS(ctx, cl, r.LogV2)

	ilToken, err := github.GenerateInstallationToken(r.LogV2, watcherServices.ILGitClient, env.Config.GitILRepositoryOwner)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating installation token for git organization [%s]", env.Config.GitILRepositoryOwner)
	}
	ilService, err := il.NewService(ctx, ilToken, r.LogV2)
	if err != nil {
		return nil, errors.Wrap(err, "error getting environment from k8s cache")
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
		return nil, errors.Wrap(err, "error instantiating git client")
	}
	fs := file.NewOSFileService()

	return &EnvironmentServices{
		StateManagerClient:     stateManagerClient,
		CloudKnitServiceClient: cloudKnitServiceClient,
		EventService:           eventServiceClient,
		ArgocdClient:           argocdClient,
		ArgoWorkflowClient:     argoworkflowClient,
		WatcherServices:        watcherServices,
		SecretsClient:          secretsClient,
		K8sClient:              k8sClient,
		ILService:              ilService,
		CompanyGitClient:       companyGitClient,
		FileService:            fs,
	}, nil
}
