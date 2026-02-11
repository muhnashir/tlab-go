// Package mariadb
package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

const (
	NOT_FOUND_ERROR = `sql: no rows in result set`
)

type Config struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	Timeout      time.Duration
	Charset      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	Type         string
	Debug        bool
}

type Meta struct {
	Total    uint64 `db:"total" json:"total"`
	Page     uint64 `json:"page"`
	Limit    uint64 `json:"limit"`
	PageNext uint64 `json:"page_next"`
}

type Adapter interface {
	Builder() *goqu.Database
	Meta(ctx context.Context, builder *goqu.SelectDataset, limit, page uint64) *Meta
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRows(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	Fetch(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	FetchRow(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping(ctx context.Context) error
	HealthCheck() error
	Error(err error) error
	Delete(ctx context.Context, table, field string, value interface{}) error
	Upsert(ctx context.Context, table interface{}, params map[string]interface{}) (uint64, error)
	UpsertWithTx(ctx context.Context, tx *sql.Tx, table interface{}, params map[string]interface{}) (uint64, error)
	SofDelete(ctx context.Context, table string, id interface{}) error
	Table(name string, stmt bool) *goqu.SelectDataset
	InsertWithTx(ctx context.Context, tx *sql.Tx, table interface{}, params map[string]interface{}) (uint64, error)
	Insert(ctx context.Context, table interface{}, params map[string]interface{}) (uint64, error)
	InsertBulk(ctx context.Context, table interface{}, params []map[string]interface{}) error
	InsertBulkWithTx(ctx context.Context, tx *sql.Tx, table interface{}, params []map[string]interface{}) error
	Update(ctx context.Context, table interface{}, field string, value interface{}, params map[string]interface{}) (uint64, error)
}
