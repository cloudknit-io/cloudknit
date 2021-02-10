package il

type config struct {
	TeamDirectory          string
	ConfigWatcherDirectory string
}

var Config = config{
	TeamDirectory:          "team",
	ConfigWatcherDirectory: "config-watcher",
}
