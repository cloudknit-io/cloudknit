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

var repoDir = "repos/zl-il"

var Config = config{
	TeamDirectory:          repoDir + "/team",
	CompanyDirectory:       repoDir + "/company",
	ConfigWatcherDirectory: repoDir + "/config-watcher",
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
