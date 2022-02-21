package il

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	git2 "github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
)

type Service struct {
	TFILGitAPI   git2.API
	TFILTempDir  string
	TFILCleanupF git2.CleanupFunc
	ZLILGitAPI   git2.API
	ZLILTempDir  string
	ZLILCleanupF git2.CleanupFunc
}

func NewService(ctx context.Context, token string) (*Service, error) {
	zlILGitAPI, err := git2.NewGoGit(ctx, token)
	if err != nil {
		return nil, err
	}

	tfILGitAPI, err := git2.NewGoGit(ctx, token)
	if err != nil {
		return nil, err
	}

	// temp clone IL repo
	tempZLILRepoDir, zlILCleanup, err := git2.CloneTemp(zlILGitAPI, env.Config.ILZLifecycleRepositoryURL)
	if err != nil {
		return nil, err
	}

	tempTFILRepoDir, tfILCleanup, err := git2.CloneTemp(tfILGitAPI, env.Config.ILTerraformRepositoryURL)
	if err != nil {
		zlILCleanup()
		return nil, err
	}

	return &Service{
		TFILGitAPI:   tfILGitAPI,
		TFILTempDir:  tempTFILRepoDir,
		TFILCleanupF: tfILCleanup,
		ZLILGitAPI:   zlILGitAPI,
		ZLILTempDir:  tempZLILRepoDir,
		ZLILCleanupF: zlILCleanup,
	}, nil
}

type config struct {
	TeamDirectory          string
	CompanyDirectory       string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamDirectory:          env.Config.ILTeamFolder,
	CompanyDirectory:       env.Config.ILCompanyFolder,
	ConfigWatcherDirectory: env.Config.ILConfigWatcherFolder,
}

func ConfigWatcherDirectoryAbsolutePath(dir string) string {
	return filepath.Join(dir, Config.ConfigWatcherDirectory)
}

func CompanyDirectoryAbsolutePath(dir string) string {
	return filepath.Join(dir, Config.CompanyDirectory)
}

func TeamDirectoryAbsolutePath(dir string) string {
	return filepath.Join(dir, Config.TeamDirectory)
}

func environmentDirectoryName(team string) string {
	return team + "-team-environment"
}

func EnvironmentDirectoryPath(team string) string {
	return Config.TeamDirectory + "/" + environmentDirectoryName(team)
}

func EnvironmentDirectoryAbsolutePath(dir string, team string) string {
	return filepath.Join(dir, EnvironmentDirectoryPath(team))
}

func environmentComponentsDirectoryName(environment string) string {
	return environment + "-environment-component"
}

func EnvironmentComponentsDirectoryPath(team string, environment string) string {
	return filepath.Join(EnvironmentDirectoryPath(team), environmentComponentsDirectoryName(environment))
}

func EnvironmentComponentsDirectoryAbsolutePath(dir string, team string, environment string) string {
	return filepath.Join(dir, EnvironmentComponentsDirectoryPath(team, environment))
}

func EnvironmentComponentDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentsDirectoryPath(team, environment), component)
}

func EnvironmentComponentDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentDirectoryPath(team, environment, component))
}

func EnvironmentComponentTerraformDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentDirectoryPath(team, environment, component), "terraform")
}

func EnvironmentComponentTerraformDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentDirectoryPath(team, environment, component), "terraform")
}

func EnvironmentComponentModuleSource(moduleSource string, moduleName string) string {
	if moduleSource == "aws" {
		return fmt.Sprintf("git@github.com:terraform-aws-modules/terraform-aws-%s.git", moduleName)
	}
	return moduleSource
}

func EnvironmentComponentModulePath(modulePath string) string {
	if modulePath == "" {
		return "."
	}
	return modulePath
}
