package repository

import (
	"strings"

	"github.com/ckshitij/notify-srv/internal/notification"
	"github.com/ckshitij/notify-srv/internal/template"
)

const (
	CreateTemplateQuery = `
		INSERT INTO templates
			(name, description, channel, type, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	CreateNotificaionQuery = `
		INSERT INTO notifications
		(channel, template_version_id, recipient, payload, status, scheduled_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	GetNotificationByIDQuery = `
		SELECT
			id, channel, template_version_id,
			recipient, payload, status,
			scheduled_at, sent_at,
			created_at, updated_at
		FROM notifications
		WHERE id = ?
		LIMIT 1
	`

	GetTemplateQuery = `
		SELECT
			id,
			name,
			description,
			channel,
			type,
			active_version,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM templates
		WHERE name = ?
		  AND type = ?
		  AND channel = ?
		LIMIT 1
	`

	GetActiveVersionQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		  AND is_active = TRUE
		LIMIT 1
	`

	GetVersionQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		  AND version = ?
		LIMIT 1
	`

	ListVersionsQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		ORDER BY version DESC
	`

	FindDueNotificationQuery = `
		SELECT id
		FROM notifications
		WHERE status = ?
		  AND scheduled_at <= NOW()
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

func buildGetAllTemplatesQuery(filter template.ListTemplatesFilter) (string, []any) {
	query := `
		SELECT
			t.id,
			t.name,
			t.description,
			t.channel,
			t.type,
			t.active_version,
			v.subject,
			v.body,
			t.created_at,
			t.updated_at
		FROM templates t
		JOIN template_versions v
			ON v.template_id = t.id
		   AND v.version = t.active_version
		WHERE 1=1
	`

	args := []any{}

	if filter.Channel != nil {
		query += " AND t.channel = ?"
		args = append(args, *filter.Channel)
	}

	if filter.Type != nil {
		query += " AND t.type = ?"
		args = append(args, *filter.Type)
	}

	query += " ORDER BY t.created_at DESC"

	return query, args
}

func buildListNotificationsQuery(filter notification.NotificationFilter) (string, []any) {
	query := `SELECT id, channel, template_version_id, recipient, payload, status, scheduled_at, sent_at, created_at, updated_at FROM notifications`
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
