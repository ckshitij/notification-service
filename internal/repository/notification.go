package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/notification"
)

type NotificationRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewNotificationRepository(db *sql.DB, log logger.Logger) *NotificationRepository {
	return &NotificationRepository{db, log}
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) (int64, error) {

	recipient, _ := json.Marshal(n.Recipient)
	payload, _ := json.Marshal(n.TemplateKeyValue)

	res, err := r.db.ExecContext(ctx, CreateNotificaionQuery,
		n.Channel,
		n.TemplateVersionID,
		recipient,
		payload,
		n.Status,
		n.ScheduledAt,
	)
	if err != nil {
		r.log.Error(ctx, "failed to create notification", logger.Error(err))
		return -1, err
	}

	id, _ := res.LastInsertId()
	n.ID = id
	return id, nil
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
		r.log.Error(ctx, "failed to get notification by id ", logger.Int64("notification_id", id), logger.Error(err))
		return nil, err
	}

	json.Unmarshal(recipient, &n.Recipient)
	json.Unmarshal(payload, &n.TemplateKeyValue)

	return &n, nil
}

func (r *NotificationRepository) FindDue(ctx context.Context, limit int) ([]int64, error) {

	rows, err := r.db.QueryContext(ctx, FindDueNotificationQuery, notification.StatusScheduled, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	r.log.Debug(ctx, "get notification scheduled ids ", logger.Int("due_ids_count", len(ids)))

	return ids, nil
}

func (r *NotificationRepository) FindStuckSending(ctx context.Context, olderThan time.Duration, limit int) ([]int64, error) {

	rows, err := r.db.QueryContext(ctx, FindStuckSendingNotificationQuery, notification.StatusSending, int(olderThan.Seconds()), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	r.log.Debug(ctx, "get notification stuck ids ", logger.Int("stuck_ids_count", len(ids)))
	return ids, nil
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
			r.log.Error(ctx, "failed to get notifications ", logger.Error(err))
			return nil, err
		}

		json.Unmarshal(recipient, &n.Recipient)
		json.Unmarshal(payload, &n.TemplateKeyValue)

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

	res, err := r.db.ExecContext(ctx, `UPDATE notifications SET status = ? WHERE id = ? AND status IN (?, ?)`,
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
