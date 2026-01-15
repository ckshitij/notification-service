package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/pkg/template"
	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/go-sql-driver/mysql"
)

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

type MySQLTemplate struct {
	db  *sql.DB
	log logger.Logger
}

func NewTemplateRepository(db *sql.DB, log logger.Logger) *MySQLTemplate {
	return &MySQLTemplate{db, log}
}

func (r *MySQLTemplate) Create(ctx context.Context, tpl template.Template) (int64, error) {

	query, args := CreateTemplateQuery, []any{tpl.Name,
		tpl.Description,
		tpl.Channel,
		tpl.Type,
		tpl.Subject,
		tpl.Body,
		tpl.CreatedBy,
		tpl.UpdatedBy,
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if isDuplicateKey(err) {
			return -1, shared.ErrDuplicateTemplateRecord
		}
		r.log.Error(ctx, "failed to create template", logger.String("query", query), logger.Field{Key: "args", Value: args}, logger.Error(err))
		return -1, err
	}

	return result.LastInsertId()

}

func (r *MySQLTemplate) GetByID(ctx context.Context, templateID int64) (*template.Template, error) {

	row := r.db.QueryRowContext(ctx, GetTemplateByIDQuery, templateID)

	var t template.Template
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Channel,
		&t.Type,
		&t.IsActive,
		&t.Subject,
		&t.Body,
		&t.CreatedBy,
		&t.UpdatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		r.log.Error(ctx, "failed to get template", logger.Int64("templateID", templateID), logger.Error(err))
		return nil, err
	}

	return &t, nil
}

func (r *MySQLTemplate) List(ctx context.Context, filter template.TemplateFilter) ([]*template.Template, error) {

	query, args := buildGetAllTemplatesQuery(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.Error(ctx, "failed to list templates", logger.String("query", query), logger.Field{Key: "args", Value: args}, logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var out = []*template.Template{}
	for rows.Next() {
		var t template.Template
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Channel,
			&t.Type,
			&t.IsActive,
			&t.Subject,
			&t.Body,
			&t.CreatedBy,
			&t.UpdatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			r.log.Error(ctx, "failed to scan list templates", logger.Error(err))
			return nil, err
		}
		out = append(out, &t)
	}

	return out, nil
}
