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
