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

func TeamDirectoryName(team string) string {
	return team + "-team-environment"
}

func TeamDirectoryPath(team string) string {
	return Config.TeamDirectory + "/" + TeamDirectoryName(team)
}

func EnvironmentDirectoryName(environment string) string {
	return environment + "-environment-component"
}

func EnvironmentDirectoryPath(team string, environment string) string {
	return TeamDirectoryPath(team) + "/" + EnvironmentDirectoryName(environment)
}

func EnvironmentComponentDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentDirectoryPath(team, environment), component)
}

func EnvironmentComponentTerraformDirectoryPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentComponentDirectoryPath(team, environment, component), "terraform")
}

func TeamTFDirectory(team string) string {
	return Config.TeamTFDirectory + "/" + team + "-team-environment"
}

func EnvironmentTFDirectory(team string, environment string) string {
	return TeamDirectoryName(team) + "/" + environment + "-environment-component"
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
