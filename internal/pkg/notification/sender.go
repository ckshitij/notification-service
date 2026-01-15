package notification

import (
	"context"

	"github.com/ckshitij/notify-srv/internal/pkg/renderer"
)

type Sender interface {
	Send(ctx context.Context, n Notification, content renderer.RenderedTemplate) error
}
