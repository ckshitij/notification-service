package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
	mysqlwrapper "github.com/ckshitij/notify-srv/internal/mysql"
	"github.com/ckshitij/notify-srv/internal/pkg/template"
	"github.com/ckshitij/notify-srv/internal/shared"
	driver "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

const (
	templateCacheByID = "template:id:%d"
	cacheExpiry       = 5 * time.Minute
)

func isDuplicateKey(err error) bool {
	var mysqlErr *driver.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

type mysqlRepo struct {
	db  *mysqlwrapper.DB
	rdb *redis.Client
	log logger.Logger
}

func NewTemplateRepository(db *mysqlwrapper.DB, rdb *redis.Client, log logger.Logger) template.TemplateRepository {
	return &mysqlRepo{db, rdb, log}
}

func (r *mysqlRepo) Create(ctx context.Context, tpl template.Template) (int64, error) {
	query, args := CreateTemplateQuery, []any{tpl.Name, tpl.Description, tpl.Channel, tpl.Type, tpl.Subject, tpl.Body, tpl.CreatedBy, tpl.UpdatedBy}
	result, err := r.db.ExecContext(ctx, "CreateNotification", query, args...)
	if err != nil {
		if isDuplicateKey(err) {
			return -1, shared.ErrDuplicateTemplateRecord
		}
		r.log.Error(ctx, "failed to create template", logger.String("query", query), logger.Field{Key: "args", Value: args}, logger.Error(err))
		return -1, err
	}

	return result.LastInsertId()
}

func (r *mysqlRepo) GetByID(ctx context.Context, templateID int64) (*template.Template, error) {
	// Check cache first
	key := fmt.Sprintf(templateCacheByID, templateID)
	cached, err := r.rdb.Get(ctx, key).Result()
	if err == nil {
		var t template.Template
		if err := json.Unmarshal([]byte(cached), &t); err == nil {
			return &t, nil
		}
	}

	row := r.db.QueryRowContext(ctx, "GetNotificationByID", GetTemplateByIDQuery, templateID)

	var t template.Template
	err = row.Scan(
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
		r.log.Info(ctx, "notification not found", logger.Int64("templateID", templateID))
		return nil, shared.ErrRecordNotFound
	}

	if err != nil {
		r.log.Error(ctx, "failed to get template", logger.Int64("templateID", templateID), logger.Error(err))
		return nil, err
	}

	// Cache the result
	serialized, _ := json.Marshal(t)
	r.rdb.Set(ctx, key, serialized, cacheExpiry)

	return &t, nil
}

func (r *mysqlRepo) CacheReloadSystemTemplates(ctx context.Context) error {
	sysTempl := shared.TemplateType(shared.SystemTemplate)
	filter := template.TemplateFilter{
		Type: &sysTempl,
	}

	templates, err := r.List(ctx, filter)
	if err != nil {
		r.log.Error(ctx, "failed to list system templates", logger.Error(err))
		return err
	}

	for _, t := range templates {
		serialized, err := json.Marshal(t)
		if err != nil {
			r.log.Warn(ctx, "failed to serialize template",
				logger.Int64("templateID", t.ID),
				logger.Error(err),
			)
			continue
		}

		pipe := r.rdb.Pipeline()
		pipe.Set(ctx, fmt.Sprintf(templateCacheByID, t.ID), serialized, cacheExpiry)

		if _, err := pipe.Exec(ctx); err != nil {
			r.log.Warn(ctx, "failed to reload template cache",
				logger.Int64("templateID", t.ID),
				logger.Error(err),
			)
		}
	}

	r.log.Info(ctx, "system templates cache reloaded",
		logger.Int("count", len(templates)),
	)

	return nil
}

func (r *mysqlRepo) List(ctx context.Context, filter template.TemplateFilter) ([]*template.Template, error) {
	query, args := buildGetAllTemplatesQuery(filter)
	rows, err := r.db.QueryContext(ctx, "ListNotification", query, args...)
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

func (r *mysqlRepo) InvalidateTemplateCache(ctx context.Context, templateID int64) error {
	keys := []string{
		fmt.Sprintf(templateCacheByID, templateID),
	}

	if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
		r.log.Error(ctx, "failed to invalidate template cache",
			logger.Field{Key: "keys", Value: keys},
			logger.Error(err),
		)
		return err
	}
	return nil
}
