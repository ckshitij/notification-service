package store

import (
	"strings"

	"github.com/ckshitij/notify-srv/internal/pkg/notification"
)

const (
	CreateNotificaionQuery = `
		INSERT INTO notifications
		(channel, template_id, recipient, template_kv, status, scheduled_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	GetNotificationByIDQuery = `
		SELECT
			id, channel, template_id,
			recipient, template_kv, status,
			scheduled_at, sent_at,
			created_at, updated_at
		FROM notifications
		WHERE id = ?
		LIMIT 1
	`

	FindDueNotificationQuery = `
		SELECT id
		FROM notifications
		WHERE status = ?
		  AND scheduled_at <= UTC_TIMESTAMP()
		ORDER BY scheduled_at
		LIMIT ?
	`

	FindStuckSendingNotificationQuery = `
		SELECT id
		FROM notifications
		WHERE status = ?
		  AND updated_at < NOW() - INTERVAL ? SECOND
		LIMIT ?
	`
)

func buildListNotificationsQuery(filter notification.NotificationFilter) (string, []any) {
	query := `SELECT id, channel, template_id, recipient, template_kv, status, scheduled_at, sent_at, created_at, updated_at FROM notifications`
	args := []any{}
	conditions := []string{}

	if filter.Channel != nil {
		conditions = append(conditions, "channel = ?")
		args = append(args, *filter.Channel)
	}
	if filter.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *filter.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}
