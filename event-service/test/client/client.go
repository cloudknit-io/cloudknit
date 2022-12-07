package main

import (
	"github.com/cloudknit-io/cloudknit/event-service/internal/zlog"
	"github.com/r3labs/sse/v2"
)

func main() {
	url := "http://localhost:8082"
	log := zlog.NewLogger()
	client := sse.NewClient(url)

	log.Infof("Connecting test client with event stream at %s", url)

	if err := client.SubscribeRaw(func(msg *sse.Event) {
		// Got some data!
		log.Infof("data: %s", string(msg.Data))
	}); err != nil {
		log.Errorf("error received from stream: %v", err)
	}
}
