package notification

import (
	"context"
	"time"
)

type Service interface {
	SendNow(ctx context.Context, n *Notification) (int64, error)
	Schedule(ctx context.Context, n *Notification, when time.Time) (int64, error)
	Process(ctx context.Context, notificationID int64) error
	GetByID(ctx context.Context, notificationID int64) (*Notification, error)
	List(ctx context.Context, filter NotificationFilter) ([]*Notification, error)
}
