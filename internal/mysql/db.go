package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/metrics"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func New(cfg config.MySQLConfig) (*DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}

	// Pooling
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connectivity (startup check)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (d *DB) ExecContext(ctx context.Context, queryName, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	res, err := d.conn.ExecContext(ctx, query, args...)
	duration := time.Since(start).Milliseconds()
	metrics.SQLQueryDuration.WithLabelValues(queryName).Observe(float64(duration))
	return res, err
}

func (d *DB) QueryContext(ctx context.Context, queryName, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	rows, err := d.conn.QueryContext(ctx, query, args...)
	duration := time.Since(start).Milliseconds()
	metrics.SQLQueryDuration.WithLabelValues(queryName).Observe(float64(duration))
	return rows, err
}

func (d *DB) QueryRowContext(ctx context.Context, queryName, query string, args ...any) *sql.Row {
	start := time.Now()
	row := d.conn.QueryRowContext(ctx, query, args...)
	duration := time.Since(start).Milliseconds()
	metrics.SQLQueryDuration.WithLabelValues(queryName).Observe(float64(duration))
	return row
}

func (d *DB) Conn() *sql.DB {
	return d.conn
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return d.conn.PingContext(ctx)
}
