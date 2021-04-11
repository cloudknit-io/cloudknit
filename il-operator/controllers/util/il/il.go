package il

import (
	"fmt"

	env "github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
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

func EnvironmentComponentDirectory(teamName string, envName string) string {
	return EnvironmentDirectory(teamName) + "/" + envName + "-environment-component"
}

func SSHKeyName() string {
	return env.Config.CompanyName + "-ssh"
}

func RepoName(companyName string) string {
	return companyName + "-il"
}

func RepoURL(owner string, companyName string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, RepoName(companyName))
}

func EnvironmentDirectory(teamName string) string {
	return Config.TeamDirectory + "/" + teamName + "-team-environment"
}

func EnvComponentModuleSource(moduleSource string, moduleName string) string {
	if moduleSource == "aws" {
		return "git@github.com:terraform-aws-modules/terraform-aws-" + moduleName + ".git"
	}
	return moduleSource
}

func EnvComponentModulePath(modulePath string) string {
	if modulePath == "" {
		return "."
	} else {
		return modulePath
	}
}
