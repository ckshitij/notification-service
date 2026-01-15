package mysql

import "github.com/ckshitij/notify-srv/internal/pkg/template"

const (
	CreateTemplateQuery = `
		INSERT INTO templates
			(name, description, channel, type, subject, body, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	GetTemplateByIDQuery = `
		SELECT
			id,
			name,
			description,
			channel,
			type,
			is_active,
			IFNULL(subject, ''), 
			body,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM templates
		WHERE id = ?
	`

	GetTemplateBySlugQuery = `
		SELECT
			id,
			name,
			description,
			channel,
			type,
			is_active,
			IFNULL(subject, ''),
			body,
			created_by,
			updated_by,
			created_at,
			updated_at
		FROM templates
		WHERE name = ?
	`
)

func buildGetAllTemplatesQuery(filter template.TemplateFilter) (string, []any) {
	query := `
		SELECT id, name, description, channel, type, is_active, IFNULL(subject, ''), 
			body, created_by, updated_by, created_at, updated_at 
		FROM templates
		WHERE 1=1
	`

	args := []any{}

	if filter.Channel != nil {
		query += " AND channel = ?"
		args = append(args, *filter.Channel)
	}

	if filter.Type != nil {
		query += " AND type = ?"
		args = append(args, *filter.Type)
	}

	if filter.IsActive != nil {
		query += " AND is_active = ?"
		args = append(args, *filter.IsActive)
	}

	if filter.Name != nil {
		query += " AND name = ?"
		args = append(args, *filter.Name)
	}

	query += " ORDER BY updated_at DESC "

	// pagination
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	return query, args
}
