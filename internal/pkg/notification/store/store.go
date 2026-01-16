package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
	mysqlwrapper "github.com/ckshitij/notify-srv/internal/mysql"
	"github.com/ckshitij/notify-srv/internal/pkg/notification"
	"github.com/ckshitij/notify-srv/internal/shared"
	driver "github.com/go-sql-driver/mysql"
)

func isFKViolation(err error) bool {
	var mysqlErr *driver.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1452
	}
	return false
}

type notificationStore struct {
	db  *mysqlwrapper.DB
	log logger.Logger
}

func NewNotificationRepository(db *mysqlwrapper.DB, log logger.Logger) *notificationStore {
	return &notificationStore{db, log}
}

func (r *notificationStore) Create(ctx context.Context, n *notification.Notification) (int64, error) {

	recipient, err := json.Marshal(n.Recipient)
	if err != nil {
		r.log.Error(ctx, "failed to marshal recipient", logger.Any("Recipient", n.Recipient), logger.Error(err))
		return -1, shared.ErrInvalidRecipient
	}

	payload, err := json.Marshal(n.TemplateKeyValue)
	if err != nil {
		r.log.Error(ctx, "failed to marshal TemplateKeyValue", logger.Any("TemplateKeyValue", n.Recipient), logger.Error(err))
		return -1, shared.ErrInvalidTemplateKeyValue
	}

	res, err := r.db.ExecContext(ctx, "CreateNotification", CreateNotificaionQuery,
		n.Channel,
		n.TemplateID,
		recipient,
		payload,
		n.Status,
		n.ScheduledAt,
	)
	if err != nil {
		if isFKViolation(err) {
			return -1, shared.ErrTemplateNotFound
		}
		r.log.Error(ctx, "failed to create notification", logger.Error(err))
		return -1, err
	}

	id, _ := res.LastInsertId()
	n.ID = id
	return id, nil
}

func (r *notificationStore) GetByID(ctx context.Context, id int64) (*notification.Notification, error) {
	row := r.db.QueryRowContext(ctx, "GetNotificationByID", GetNotificationByIDQuery, id)

	var (
		n         notification.Notification
		recipient []byte
		payload   []byte
	)

	err := row.Scan(
		&n.ID,
		&n.Channel,
		&n.TemplateID,
		&recipient,
		&payload,
		&n.Status,
		&n.ScheduledAt,
		&n.SentAt,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Info(ctx, "notification not found", logger.Int64("notification_id", id))
			return nil, shared.ErrRecordNotFound
		}

		r.log.Error(ctx, "failed to get notification by id", logger.Int64("notification_id", id), logger.Error(err))
		return nil, err
	}

	if err := json.Unmarshal(recipient, &n.Recipient); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(payload, &n.TemplateKeyValue); err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *notificationStore) FindDue(ctx context.Context, limit int) ([]int64, error) {

	rows, err := r.db.QueryContext(ctx, "FindDueNotifications", FindDueNotificationQuery, notification.StatusScheduled, limit)
	if err != nil {
		r.log.Error(ctx, "failed to find due notifications ", logger.Int("limit", limit), logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			r.log.Error(ctx, "failed to scan due notifications ids", logger.Error(err))
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *notificationStore) FindStuckSending(ctx context.Context, olderThan time.Duration, limit int) ([]int64, error) {

	rows, err := r.db.QueryContext(ctx, "FindStuckSendingNotifications", FindStuckSendingNotificationQuery, notification.StatusSending, int(olderThan.Seconds()), limit)
	if err != nil {
		r.log.Error(ctx, "failed to find stuck notifications ", logger.Any("olderThan", olderThan), logger.Int("limit", limit), logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			r.log.Error(ctx, "failed to scan stuck notifications ids", logger.Error(err))
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *notificationStore) List(ctx context.Context, filter notification.NotificationFilter) ([]*notification.Notification, error) {

	query, args := buildListNotificationsQuery(filter)

	rows, err := r.db.QueryContext(ctx, "ListNotifications", query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications = []*notification.Notification{}

	for rows.Next() {
		var (
			n         notification.Notification
			recipient []byte
			payload   []byte
		)

		if err := rows.Scan(
			&n.ID,
			&n.Channel,
			&n.TemplateID,
			&recipient,
			&payload,
			&n.Status,
			&n.ScheduledAt,
			&n.SentAt,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			r.log.Error(ctx, "failed to scan notifications list", logger.Error(err))
			return nil, err
		}

		json.Unmarshal(recipient, &n.Recipient)
		json.Unmarshal(payload, &n.TemplateKeyValue)

		notifications = append(notifications, &n)
	}

	if err := rows.Err(); err != nil {
		r.log.Error(ctx, "failed to find scan notifications list", logger.Error(err))
		return nil, err
	}

	return notifications, nil
}

func (r *notificationStore) MarkSent(ctx context.Context, id int64, sentAt time.Time) error {
	_, err := r.db.ExecContext(ctx, "MarkNotificationSent", `UPDATE notifications SET status = ?, sent_at = ? WHERE id = ?`, notification.StatusSent, sentAt, id)
	if err != nil {
		r.log.Error(ctx, "failed to mark notification ", logger.String("status", "sent"), logger.Int64("notificationID", id), logger.Error(err))
	}
	return err
}

func (r *notificationStore) UpdateStatus(ctx context.Context, id int64, status notification.NotificationStatus) error {
	_, err := r.db.ExecContext(ctx, "UpdateNotificationStatus", `UPDATE notifications SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		r.log.Error(ctx, "failed to update notification status ", logger.Int64("notificationID", id), logger.Error(err))
	}
	return err
}

func (r *notificationStore) AcquireForSending(ctx context.Context, id int64) (bool, error) {

	res, err := r.db.ExecContext(ctx, "AcquireNotificationForSending", `UPDATE notifications SET status = ? WHERE id = ? AND status IN (?, ?)`,
		notification.StatusSending,
		id,
		notification.StatusPending,
		notification.StatusScheduled,
	)
	if err != nil {
		r.log.Error(ctx, "failed to acquire notification ", logger.String("status", "sending"), logger.Int64("notificationID", id), logger.Error(err))
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows == 1, nil
}
