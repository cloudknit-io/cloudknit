package il

import (
	"fmt"
	"path/filepath"
)

type config struct {
	TeamTFDirectory        string
	TeamDirectory          string
	CompanyDirectory       string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamTFDirectory:        "team-tf",
	TeamDirectory:          "team",
	CompanyDirectory:       "company",
	ConfigWatcherDirectory: "config-watcher",
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

func EnvironmentDirectoryName(team string) string {
	return team + "-team-environment"
}

func EnvironmentDirectoryPath(team string) string {
	return Config.TeamDirectory + "/" + EnvironmentDirectoryName(team)
}

func EnvironmentDirectoryAbsolutePath(dir string, team string) string {
	return filepath.Join(dir, EnvironmentDirectoryPath(team))
}

func EnvironmentComponentsDirectoryName(environment string) string {
	return environment + "-environment-component"
}

func EnvironmentComponentsDirectoryPath(team string, environment string) string {
	return filepath.Join(EnvironmentDirectoryPath(team), EnvironmentComponentsDirectoryName(environment))
}

func EnvironmentComponentsDirectoryAbsolutePath(dir string, team string, environment string) string {
	return filepath.Join(dir, EnvironmentComponentsDirectoryPath(team, environment))
}

func EnvironmentComponentDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentsDirectoryPath(team, environment), component)
}

func EnvironmentComponentDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentsDirectoryPath(team, environment), component)
}

func EnvironmentComponentTerraformDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentDirectoryPath(team, environment, component), "terraform")
}

func EnvironmentComponentTerraformDirectoryAbsolutePath(dir string, team string, environment string, component string) string {
	return filepath.Join(dir, EnvironmentComponentDirectoryPath(team, environment, component), "terraform")
}

func TeamTFDirectory(team string) string {
	return Config.TeamTFDirectory + "/" + team + "-team-environment"
}

func EnvironmentTFDirectory(team string, environment string) string {
	return EnvironmentDirectoryName(team) + "/" + environment + "-environment-component"
}

func RepoName(companyName string) string {
	return companyName + "-il"
}

func RepoURL(owner string, companyName string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, RepoName(companyName))
}

func EnvComponentModuleSource(moduleSource string, moduleName string) string {
	if moduleSource == "aws" {
		return fmt.Sprintf("git@github.com:terraform-aws-modules/terraform-aws-%s.git", moduleName)
	}
	return moduleSource
}

func EnvComponentModulePath(modulePath string) string {
	if modulePath == "" {
		return "."
	}
	return modulePath
}
