package notification

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, n *Notification) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status NotificationStatus) error
	GetByID(ctx context.Context, id int64) (*Notification, error)
	List(ctx context.Context, filter NotificationFilter) ([]*Notification, error)
	MarkSent(ctx context.Context, id int64, sentAt time.Time) error
	MarkFailed(ctx context.Context, id int64, code string, message string, metadata map[string]any) error
	AcquireForSending(ctx context.Context, id int64) (bool, error)
	FindDue(ctx context.Context, limit int) ([]NotificationScheduled, error)
	FindStuckSending(ctx context.Context, olderThan time.Duration, limit int) ([]NotificationScheduled, error)
}
