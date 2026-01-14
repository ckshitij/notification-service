package repository

import (
	"context"
	"database/sql"

	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/ckshitij/notify-srv/internal/template"
)

type TemplateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) CreateTemplate(
	ctx context.Context,
	tpl template.Template,
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

func (r *TemplateRepository) GetTemplate(ctx context.Context, name string, tplType shared.TemplateType, channel shared.Channel) (*template.Template, error) {

	row := r.db.QueryRowContext(ctx, GetTemplateQuery, name, tplType, channel)

	var tpl template.Template
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

func (r *TemplateRepository) CreateVersion(
	ctx context.Context,
	version template.TemplateVersion,
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

func (r *TemplateRepository) GetActiveVersion(
	ctx context.Context,
	templateID int64,
) (*template.TemplateVersion, error) {

	row := r.db.QueryRowContext(ctx, GetActiveVersionQuery, templateID)

	var v template.TemplateVersion
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

func (r *TemplateRepository) GetVersion(
	ctx context.Context,
	templateID int64,
	version int,
) (*template.TemplateVersion, error) {

	row := r.db.QueryRowContext(ctx, GetActiveVersionQuery, templateID, version)

	var v template.TemplateVersion
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

func (r *TemplateRepository) ListVersions(ctx context.Context, templateID int64) ([]template.TemplateVersion, error) {
	rows, err := r.db.QueryContext(ctx, ListVersionsQuery, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []template.TemplateVersion

	for rows.Next() {
		var v template.TemplateVersion
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

func (r *TemplateRepository) ListTemplatesWithActiveVersion(ctx context.Context, filter template.ListTemplatesFilter) ([]template.TemplateWithActiveVersion, error) {

	query, args := GetAllTemplatesQuery(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []template.TemplateWithActiveVersion

	for rows.Next() {
		var t template.TemplateWithActiveVersion
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
