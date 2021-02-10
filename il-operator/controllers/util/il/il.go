package il

type config struct {
	TeamDirectory          string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamDirectory:          "teams",
	ConfigWatcherDirectory: "config_watchers",
}
