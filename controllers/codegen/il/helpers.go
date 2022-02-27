package il

import (
	"fmt"
	"path/filepath"
)

const (
	terraformDir  = "terraform"
	argocdAppsDir = "argocd"
)

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
	return filepath.Join(EnvironmentComponentDirectoryPath(team, environment, component), terraformDir)
}

func EnvironmentComponentTerraformDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentTerraformDirectoryPath(team, environment, component))
}

func EnvironmentComponentArgocdAppsDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentsDirectoryPath(team, environment), argocdAppsDir, component)
}

func EnvironmentComponentArgocdAppsDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentArgocdAppsDirectoryPath(team, environment, component))
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
