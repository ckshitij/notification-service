package renderer

import (
	"bytes"
	"fmt"
	"text/template"
)

type RenderedTemplate struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Renderer interface {
	Render(subject, body string, data map[string]any) (RenderedTemplate, error)
}

type GoTemplateRenderer struct {
}

func NewGoTemplateRenderer() Renderer {
	return &GoTemplateRenderer{}
}

func (r *GoTemplateRenderer) Render(subject, body string, data map[string]any) (RenderedTemplate, error) {

	var result RenderedTemplate

	// Render subject (if present)
	if subject != "" {
		subject, err := renderString(subject, data)
		if err != nil {
			return result, fmt.Errorf("render subject: %w", err)
		}
		result.Subject = subject
	}

	// Render body (required)
	body, err := renderString(body, data)
	if err != nil {
		return result, fmt.Errorf("render body: %w", err)
	}
	result.Body = body

	return result, nil
}

func renderString(tpl string, data map[string]any) (string, error) {

	t, err := template.New("tpl").Option("missingkey=error").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
