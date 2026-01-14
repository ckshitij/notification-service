package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ckshitij/notify-srv/internal/notification"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) error {

	recipient, _ := json.Marshal(n.Recipient)
	payload, _ := json.Marshal(n.Payload)

	res, err := r.db.ExecContext(ctx, CreateNotificaionQuery,
		n.Channel,
		n.TemplateVersionID,
		recipient,
		payload,
		n.Status,
		n.ScheduledAt,
	)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	n.ID = id
	return nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id int64, status notification.NotificationStatus) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *NotificationRepository) GetByID(ctx context.Context, id int64) (*notification.Notification, error) {

	row := r.db.QueryRowContext(ctx, GetNotificationByIDQuery, id)

	var (
		n         notification.Notification
		recipient []byte
		payload   []byte
	)

	if err := row.Scan(
		&n.ID,
		&n.Channel,
		&n.TemplateVersionID,
		&recipient,
		&payload,
		&n.Status,
		&n.ScheduledAt,
		&n.SentAt,
		&n.CreatedAt,
		&n.UpdatedAt,
	); err != nil {
		return nil, err
	}

	json.Unmarshal(recipient, &n.Recipient)
	json.Unmarshal(payload, &n.Payload)

	return &n, nil
}

func (r *NotificationRepository) List(ctx context.Context, filter notification.NotificationFilter) ([]*notification.Notification, error) {

	query, args := buildListNotificationsQuery(filter)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*notification.Notification

	for rows.Next() {
		var (
			n         notification.Notification
			recipient []byte
			payload   []byte
		)

		if err := rows.Scan(
			&n.ID,
			&n.Channel,
			&n.TemplateVersionID,
			&recipient,
			&payload,
			&n.Status,
			&n.ScheduledAt,
			&n.SentAt,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}

		json.Unmarshal(recipient, &n.Recipient)
		json.Unmarshal(payload, &n.Payload)

		notifications = append(notifications, &n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkSent(ctx context.Context, id int64, sentAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET status = ?, sent_at = ? WHERE id = ?`, notification.StatusSent, sentAt, id)
	return err
}

func (r *NotificationRepository) AcquireForSending(ctx context.Context, id int64) (bool, error) {

	res, err := r.db.ExecContext(ctx, `
		UPDATE notifications
		SET status = ?
		WHERE id = ?
		  AND status IN (?, ?)
	`,
		notification.StatusSending,
		id,
		notification.StatusPending,
		notification.StatusScheduled,
	)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
}
