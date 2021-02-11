package il

type config struct {
	TeamDirectory          string
	ConfigWatcherDirectory string
	CompanyConfigRepo      string
}

var Config = config{
	TeamDirectory:          "team",
	ConfigWatcherDirectory: "config-watcher",
	CompanyConfigRepo:      "git@github.com:CompuZest/zmart-config.git",
}

func EnvironmentComponentDirectory(teamName string, envName string) string {
	return EnvironmentDirectory(teamName) + "/" + envName + "-environment-component"
}

func EnvironmentDirectory(teamName string) string {
	return Config.TeamDirectory + "/" + teamName + "-team-environment"
}
