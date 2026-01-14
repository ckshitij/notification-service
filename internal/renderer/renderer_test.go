package renderer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTemplate(t *testing.T) {
	r := NewGoTemplateRenderer()

	var subject = "Hello {{.Name}}"
	var body = "Welcome to {{.App}}"

	out, err := r.Render(subject, body, map[string]any{
		"Name": "User",
		"App":  "NotifyX",
	})

	require.NoError(t, err)
	require.Equal(t, "Hello User", out.Subject)
	require.Equal(t, "Welcome to NotifyX", out.Body)
}
