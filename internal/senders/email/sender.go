package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/ckshitij/notify-srv/internal/notification"
	"github.com/ckshitij/notify-srv/internal/renderer"
)

type Sender struct {
	host string
	port int
	from string
	auth smtp.Auth
}

func New(
	host string,
	port int,
	from string,
	username string,
	password string,
) *Sender {

	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	return &Sender{
		host: host,
		port: port,
		from: from,
		auth: auth,
	}
}

func (s *Sender) Send(
	ctx context.Context,
	n notification.Notification,
	content renderer.RenderedTemplate,
) error {

	if n.Recipient.Email == nil {
		return fmt.Errorf("email recipient missing")
	}

	msg := fmt.Appendf(nil,
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		*n.Recipient.Email,
		content.Subject,
		content.Body,
	)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	return smtp.SendMail(
		addr,
		s.auth,
		s.from,
		[]string{*n.Recipient.Email},
		msg,
	)
}
