package notifier

import "context"

type API interface {
	Notify(ctx context.Context, n *Notification) error
}
