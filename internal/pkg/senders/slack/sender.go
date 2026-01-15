package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ckshitij/notify-srv/internal/pkg/notification"
	"github.com/ckshitij/notify-srv/internal/pkg/renderer"
)

type Sender struct {
	webhookURL string
	client     *http.Client
}

func New(webhookURL string) *Sender {
	return &Sender{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (s *Sender) Send(
	ctx context.Context,
	n notification.Notification,
	content renderer.RenderedTemplate,
) error {

	payload := map[string]string{
		"text": content.Body,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.webhookURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook failed: %s %v", resp.Status, err)
	}

	return nil
}
