package template

import "github.com/ckshitij/notify-srv/internal/shared"

type CreateTemplateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Channel     shared.Channel `json:"channel"`
}

type AddVersionRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type RenderRequest struct {
	TemplateKeyValue map[string]any `json:"template_key_value"`
}
