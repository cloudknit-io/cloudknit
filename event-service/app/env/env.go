package env

import "os"

type Cfg struct {
	App            string
	DevMode        string
	EnableNewRelic string
	NewRelicAPIKey string
	Instance       string
	Environment    string
}

func Config() *Cfg {
	return &Cfg{
		App:            "zlifecycle-event-service",
		DevMode:        getOr("DEV_MODE", "false"),
		EnableNewRelic: getOr("ENABLE_NEW_RELIC", "false"),
		NewRelicAPIKey: os.Getenv("NEW_RELIC_API_KEY"),
		Instance:       getOr("META_INSTANCE", "dev"),
		Environment:    getOr("ENVIRONMENT", "dev"),
	}
}

func getOr(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if exists && val != "" {
		return val
	}
	return defaultValue
}
