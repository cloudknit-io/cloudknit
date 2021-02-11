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

func environmentComponentDirectory(teamName string, envName string) string {
	return environmentDirectory(teamName) + "/" + envName + "-environment-component"
}

func environmentDirectory(teamName string) string {
	return Config.TeamDirectory + "/" + teamName + "-team-environment"
}
