package template

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTemplate(t *testing.T) {
	r := NewGoTemplateRenderer()

	tpl := TemplateVersion{
		Subject: "Hello {{.Name}}",
		Body:    "Welcome to {{.App}}",
	}

	out, err := r.Render(tpl, map[string]any{
		"Name": "User",
		"App":  "NotifyX",
	})

	require.NoError(t, err)
	require.Equal(t, "Hello User", out.Subject)
	require.Equal(t, "Welcome to NotifyX", out.Body)
}
