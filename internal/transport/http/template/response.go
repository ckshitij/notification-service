package template

import "time"

type RenderResponse struct {
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body"`
}

type VersionResponse struct {
	Version   int       `json:"version"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
