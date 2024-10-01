package sqldb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/Red-Sock/toolbox/closer"
	"github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka/resources"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

func New(cfg resources.SqlResource) (DB, error) {
	dialect := cfg.SqlDialect()
	connStr := cfg.ConnectionString()

	conn, err := sql.Open(dialect, connStr)
	if err != nil {
		return nil, errors.Wrap(err, "error checking connection to postgres")
	}

	closer.Add(func() error {
		return conn.Close()
	})

	goose.SetLogger(logrus.StandardLogger())
	err = goose.SetDialect(dialect)
	if err != nil {
		return nil, errors.Wrap(err, "error setting dialect")
	}

	mig := cfg.MigrationFolder()
	if mig == "" {
		mig = "./migrations"
	}

	err = goose.Up(conn, mig)
	if err != nil {
		return nil, errors.Wrap(err, "error performing up")
	}

	return conn, nil
}

type DB interface {
	PingContext(ctx context.Context) error
	Ping() error
	Close() error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	Stats() sql.DBStats
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
	Driver() driver.Driver
	Conn(ctx context.Context) (*sql.Conn, error)
}
