package uinotifier

import (
	"bytes"
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/notifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (u *UINotifier) Notify(ctx context.Context, n *notifier.Notification) error {
	u.log.WithFields(logrus.Fields{
		"message":     n.Message,
		"messageType": n.MessageType,
	}).Info("Sending UI notification")

	notificationEndpoint := fmt.Sprintf("%s/reconciliation/api/v1/notification/save", u.apiURL)

	jsonBody, err := common.ToJSON(n)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", notificationEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := common.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request to %s: %w", notificationEndpoint, err)
	}
	defer common.CloseBody(resp.Body)

	if resp.StatusCode != 201 {
		return fmt.Errorf("bad status code received: %d", resp.StatusCode)
	}

	return nil
}

var _ notifier.Notifier = &UINotifier{}
