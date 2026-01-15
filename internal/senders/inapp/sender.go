package inapp

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ckshitij/notify-srv/internal/notification"
	"github.com/ckshitij/notify-srv/internal/renderer"
)

type Sender struct {
	db *sql.DB
}

func New(db *sql.DB) *Sender {
	return &Sender{db: db}
}

func (s *Sender) Send(
	ctx context.Context,
	n notification.Notification,
	content renderer.RenderedTemplate,
) error {

	if n.Recipient.InAppUser == nil {
		return fmt.Errorf("in_app user missing in recipient")
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO in_app_notifications
		(notification_id, user_id, body)
		VALUES (?, ?, ?)
	`,
		n.ID,
		*n.Recipient.InAppUser,
		content.Body,
	)

	return err
}
