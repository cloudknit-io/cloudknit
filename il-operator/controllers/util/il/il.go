package il

import (
	"fmt"
	"path/filepath"
)

type config struct {
	TeamDirectory          string
	CompanyDirectory       string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamDirectory:          "team",
	CompanyDirectory:       "company",
	ConfigWatcherDirectory: "config-watcher",
}

func TeamDirectory(team string) string {
	return Config.TeamDirectory + "/" + team + "-team-environment"
}

func EnvironmentDirectory(team string, environment string) string {
	return TeamDirectory(team) + "/" + environment + "-environment-component"
}

func EnvironmentComponentTerraformIlPath(team string, environment string, component string) string {
	return filepath.Join(EnvironmentDirectory(team, environment), component, "terraform")
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
