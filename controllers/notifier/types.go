package notifier

import "time"

type MessageType string

const (
	INFO  = "0"
	WARN  = "1"
	ERROR = "2"
)

type Notification struct {
	Company     string      `json:"companyId"`
	Team        string      `json:"teamName"`
	Environment string      `json:"environmentName"`
	Timestamp   time.Time   `json:"timestamp"`
	MessageType MessageType `json:"messageType"`
	Message     string      `json:"message"`
	Debug       interface{} `json:"debug"`
}
