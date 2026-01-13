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

func GetAllTemplatesQuery(filter ListTemplatesFilter) (string, []any) {
	query := `
		SELECT
			t.id,
			t.name,
			t.description,
			t.channel,
			t.type,
			t.active_version,
			v.subject,
			v.body,
			t.created_at,
			t.updated_at
		FROM templates t
		JOIN template_versions v
			ON v.template_id = t.id
		   AND v.version = t.active_version
		WHERE 1=1
	`

	args := []any{}

	if filter.Channel != nil {
		query += " AND t.channel = ?"
		args = append(args, *filter.Channel)
	}

	if filter.Type != nil {
		query += " AND t.type = ?"
		args = append(args, *filter.Type)
	}

	query += " ORDER BY t.created_at DESC"

	return query, args
}
