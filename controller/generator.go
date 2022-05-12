package controller

import (
	"context"
	"fmt"

	"github.com/compuzest/zlifecycle-il-operator/controller/external/aws/awseks"

	"github.com/compuzest/zlifecycle-il-operator/controller/external/secrets"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/tfgen/tfvar"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/zerrors"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/apps"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/overlay"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/tfgen"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/git"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	TerraformSubtypeCustom = "custom"
	TerraformSubtypeIL     = "il"
)

// COMPANY

func generateAndSaveCompanyApp(fileAPI file.API, company *stablev1.Company, ilRepoDir string) error {
	companyApp := argocd.GenerateCompanyApp(company)

	return fileAPI.SaveYamlFile(*companyApp, il.CompanyDirectoryAbsolutePath(ilRepoDir), company.Spec.CompanyName+".yaml")
}

func generateAndSaveCompanyConfigWatcher(fileAPI file.API, company *stablev1.Company, ilRepoDir string) error {
	companyConfigWatcherApp := argocd.GenerateCompanyConfigWatcherApp(
		company.Spec.CompanyName,
		company.Spec.ConfigRepo.Source,
		company.Spec.ConfigRepo.Path,
	)

	return fileAPI.SaveYamlFile(*companyConfigWatcherApp, il.ConfigWatcherDirectoryAbsolutePath(ilRepoDir), company.Spec.CompanyName+".yaml")
}

// TEAM

func generateAndSaveTeamApp(fileAPI file.API, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamApp := argocd.GenerateTeamApp(team)

	return fileAPI.SaveYamlFile(*teamApp, il.TeamDirectoryAbsolutePath(ilRepoDir), filename)
}

func generateAndSaveConfigWatchers(fileAPI file.API, team *stablev1.Team, filename string, ilRepoDir string) error {
	teamConfigWatcherApp := argocd.GenerateTeamConfigWatcherApp(team)

	return fileAPI.SaveYamlFile(*teamConfigWatcherApp, il.ConfigWatcherDirectoryAbsolutePath(ilRepoDir), filename)
}

// ENVIRONMENT

func generateAndSaveWorkflowOfWorkflows(
	fileAPI file.API,
	ilService *il.Service,
	environment *stablev1.Environment,
	tfcfg *secrets.TerraformStateConfig,
) error {
	// WIP, below command is for testing
	// experimentalworkflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	// if err := fileAPI.SaveYamlFile(*experimentalworkflow, envComponentDirectory, "/experimental_wofw.yaml"); err != nil {
	// 	return err
	// }
	ilEnvComponentDirectory := il.EnvironmentComponentsDirectoryAbsolutePath(ilService.ZLILTempDir, environment.Spec.TeamName, environment.Spec.EnvName)

	workflow := argoworkflow.GenerateLegacyWorkflowOfWorkflows(environment, tfcfg)
	return fileAPI.SaveYamlFile(*workflow, ilEnvComponentDirectory, "/wofw.yaml")
}

func generateAndSaveEnvironmentApp(fileService file.API, environment *stablev1.Environment, envDirectory string) error {
	envApp := argocd.GenerateEnvironmentApp(environment)
	envYAML := fmt.Sprintf("%s-environment.yaml", environment.Spec.EnvName)

	return fileService.SaveYamlFile(*envApp, envDirectory, envYAML)
}

func generateAndSaveEnvironmentComponents(
	ctx context.Context,
	log *logrus.Entry,
	ilService *il.Service,
	fileService file.API,
	gitReconciler gitreconciler.API,
	gitClient git.API,
	k8sClient awseks.API,
	argocdClient argocd.API,
	e *stablev1.Environment,
	tfcfg *secrets.TerraformStateConfig,
) error {
	for _, ec := range e.Spec.Components {
		ecDirectory := il.EnvironmentComponentsDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName, e.Spec.EnvName)

		application := argocd.GenerateEnvironmentComponentApps(e, ec)
		if err := fileService.SaveYamlFile(*application, ecDirectory, fmt.Sprintf("%s.yaml", ec.Name)); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error saving yaml file"))
		}

		gitReconcilerKey := kClient.ObjectKey{Name: e.Name, Namespace: e.Namespace}

		switch ec.Type {
		case "terraform":
			if err := generateTerraformComponent(
				gitReconciler,
				ilService,
				gitClient,
				fileService,
				e,
				ec,
				&gitReconcilerKey,
				tfcfg,
				log,
			); err != nil {
				return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating component terraform"))
			}
		case "argocd":
			if err := generateAppsComponent(
				ctx,
				gitReconciler,
				ilService,
				gitClient,
				fileService,
				k8sClient,
				argocdClient,
				e,
				ec,
				&gitReconcilerKey,
				log,
			); err != nil {
				return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating component apps"))
			}
		default:
			return errors.Errorf("invalid environment component type: %s", ec.Type)
		}
	}

	return nil
}

func generateTerraformComponent(
	gitReconciler gitreconciler.API,
	ilService *il.Service,
	gitClient git.API,
	fileService file.API,
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	key *kClient.ObjectKey,
	tfcfg *secrets.TerraformStateConfig,
	log *logrus.Entry,
) error {
	tfDirectory := il.EnvironmentComponentTerraformDirectoryAbsolutePath(ilService.TFILTempDir, e.Spec.TeamName, e.Spec.EnvName, ec.Name)

	log.WithFields(logrus.Fields{
		"type":      ec.Type,
		"directory": tfDirectory,
	}).Infof("Generating terraform code for environment component %s", ec.Name)

	if ec.Variables != nil {
		log.WithFields(logrus.Fields{
			"component": ec.Name,
			"type":      ec.Type,
		}).Infof("Generating tfvars file from inline variables for component %s", ec.Name)
		fileName := fmt.Sprintf("%s.tfvars", ec.Name)
		if err := tfvar.GenerateTFVarsFile(fileService, ec.Variables, tfDirectory, fileName); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error saving tfvars to file"))
		}
	}

	var generatedTFVars string
	if ec.VariablesFile != nil {
		key := kClient.ObjectKey{Namespace: e.Namespace, Name: e.Name}
		extracted, err := tfvar.GetVariablesFromTfvarsFile(gitReconciler, gitClient, log, &key, ec, e.Spec.ZLocals)
		if err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error reading variables from tfvars file"))
		}

		generatedTFVars = extracted
	}

	// Deleting terraform folder so that it gets recreated so that any dangling files are cleaned up
	if err := fileService.RemoveAll(tfDirectory); err != nil {
		return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error deleting terraform directory"))
	}

	vars := tfgen.NewTemplateVariablesFromEnvironment(e, ec, generatedTFVars, tfcfg)
	if ec.Subtype == TerraformSubtypeCustom {
		if err := tfgen.GenerateCustomTerraform(fileService, gitClient, ec.Module.Source, tfDirectory); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating custom terraform"))
		}
	} else {
		if err := tfgen.GenerateTerraform(fileService, vars, tfDirectory); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating terraform"))
		}
	}

	if err := overlay.GenerateOverlayFiles(log, fileService, gitClient, gitReconciler, key, ec, tfDirectory); err != nil {
		return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating overlay files"))
	}

	return nil
}

func generateAppsComponent(
	ctx context.Context,
	gitReconciler gitreconciler.API,
	ilService *il.Service,
	gitClient git.API,
	fileAPI file.API,
	k8sClient awseks.API,
	argocdClient argocd.API,
	e *stablev1.Environment,
	ec *stablev1.EnvironmentComponent,
	key *kClient.ObjectKey,
	log *logrus.Entry,
) error {
	info, err := argocd.RegisterNewCluster(ctx, k8sClient, argocdClient, "dev-checkout-staging-eks", log)
	if err != nil {
		return errors.Wrap(err, "error registering cluster %s for environment %s")
	}

	appDirectory := il.EnvironmentComponentArgocdAppsDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName, e.Spec.EnvName, ec.Name)

	// Deleting app folder so that it gets recreated so that any dangling files are cleaned up
	if err := fileAPI.RemoveAll(appDirectory); err != nil {
		return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error deleting application directory"))
	}

	log.WithFields(logrus.Fields{
		"type":      ec.Type,
		"directory": appDirectory,
	}).Infof("Generating argocd applications for environment component %s", ec.Name)

	if err := apps.GenerateArgocdApps(log, fileAPI, gitClient, gitReconciler, key, e, info.Endpoint, appDirectory); err != nil {
		return zerrors.NewEnvironmentComponentError(ec.Name, errors.Wrap(err, "error generating overlay files"))
	}

	return nil
}
