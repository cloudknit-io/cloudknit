package uinotifier

import (
	"bytes"
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/notifier"
	"net/http"

	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func (u *UINotifier) Notify(ctx context.Context, n *notifier.Notification) error {
	u.log.WithFields(logrus.Fields{
		"message":     n.Message,
		"messageType": n.MessageType,
	}).Info("Sending UI notification")

	notificationEndpoint := fmt.Sprintf("%s/reconciliation/api/v1/notification/save", u.apiURL)

	jsonBody, err := util.ToJSON(n)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notificationEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := util.GetHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request to %s: %w", notificationEndpoint, err)
	}
	defer util.CloseBody(resp.Body)

	if resp.StatusCode != 201 {
		return errors.Errorf("bad status code received: %d", resp.StatusCode)
	}

	return nil
}

var _ notifier.API = &UINotifier{}
