package il

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

func EnvironmentDirectory(teamName string) string {
	return Config.TeamDirectory + "/" + teamName + "-team-environment"
}
