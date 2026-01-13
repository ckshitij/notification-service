package template

const (
	CreateTemplateQuery = `
		INSERT INTO templates
			(name, description, channel, type, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	GetTemplateQuery = `
		SELECT
			id,
			name,
			description,
			channel,
			type,
			active_version,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM templates
		WHERE name = ?
		  AND type = ?
		  AND channel = ?
		LIMIT 1
	`

	GetActiveVersionQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		  AND is_active = TRUE
		LIMIT 1
	`

	GetVersionQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		  AND version = ?
		LIMIT 1
	`

	ListVersionsQuery = `
		SELECT
			id,
			template_id,
			version,
			subject,
			body,
			is_active,
			created_at
		FROM template_versions
		WHERE template_id = ?
		ORDER BY version DESC
	`
)
