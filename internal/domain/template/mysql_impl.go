package template

import (
	"context"
	"database/sql"

	"github.com/ckshitij/notification-srv/internal/domain/shared"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateTemplate(
	ctx context.Context,
	tpl Template,
) error {

	_, err := r.db.ExecContext(
		ctx,
		CreateTemplateQuery,
		tpl.Name,
		tpl.Description,
		tpl.Channel,
		tpl.Type,
		tpl.CreatedBy,
		tpl.UpdatedBy,
	)

	return err
}

func (r *MySQLRepository) GetTemplate(ctx context.Context, name string, tplType shared.TemplateType, channel shared.Channel) (*Template, error) {

	row := r.db.QueryRowContext(ctx, GetTemplateQuery, name, tplType, channel)

	var tpl Template
	err := row.Scan(
		&tpl.ID,
		&tpl.Name,
		&tpl.Description,
		&tpl.Channel,
		&tpl.Type,
		&tpl.ActiveVersion,
		&tpl.CreatedBy,
		&tpl.UpdatedBy,
		&tpl.CreatedAt,
		&tpl.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &tpl, nil
}

func (r *MySQLRepository) CreateVersion(
	ctx context.Context,
	version TemplateVersion,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deactivate old versions
	const deactivate = `
		UPDATE template_versions
		SET is_active = FALSE
		WHERE template_id = ?
	`
	if _, err := tx.ExecContext(ctx, deactivate, version.TemplateID); err != nil {
		return err
	}

	// Insert new version
	const insert = `
		INSERT INTO template_versions
		(template_id, version, subject, body, is_active)
		VALUES (?, ?, ?, ?, TRUE)
	`
	if _, err := tx.ExecContext(
		ctx,
		insert,
		version.TemplateID,
		version.Version,
		version.Subject,
		version.Body,
	); err != nil {
		return err
	}

	// Update active version pointer
	const updateTpl = `
		UPDATE templates
		SET active_version = ?
		WHERE id = ?
	`
	if _, err := tx.ExecContext(
		ctx,
		updateTpl,
		version.Version,
		version.TemplateID,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *MySQLRepository) GetActiveVersion(
	ctx context.Context,
	templateID int64,
) (*TemplateVersion, error) {

	row := r.db.QueryRowContext(ctx, GetActiveVersionQuery, templateID)

	var v TemplateVersion
	err := row.Scan(
		&v.ID,
		&v.TemplateID,
		&v.Version,
		&v.Subject,
		&v.Body,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *MySQLRepository) GetVersion(
	ctx context.Context,
	templateID int64,
	version int,
) (*TemplateVersion, error) {

	row := r.db.QueryRowContext(ctx, GetActiveVersionQuery, templateID, version)

	var v TemplateVersion
	err := row.Scan(
		&v.ID,
		&v.TemplateID,
		&v.Version,
		&v.Subject,
		&v.Body,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *MySQLRepository) ListVersions(ctx context.Context, templateID int64) ([]TemplateVersion, error) {
	rows, err := r.db.QueryContext(ctx, ListVersionsQuery, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []TemplateVersion

	for rows.Next() {
		var v TemplateVersion
		if err := rows.Scan(
			&v.ID,
			&v.TemplateID,
			&v.Version,
			&v.Subject,
			&v.Body,
			&v.IsActive,
			&v.CreatedAt,
		); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (r *MySQLRepository) ListTemplatesWithActiveVersion(ctx context.Context, filter ListTemplatesFilter) ([]TemplateWithActiveVersion, error) {

	query, args := GetAllTemplatesQuery(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TemplateWithActiveVersion

	for rows.Next() {
		var t TemplateWithActiveVersion
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Channel,
			&t.Type,
			&t.ActiveVersion,
			&t.Subject,
			&t.Body,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}
