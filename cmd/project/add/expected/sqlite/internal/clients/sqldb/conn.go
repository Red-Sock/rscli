package sqldb

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox/closer"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/pressly/goose/v3"
)

func New(cfg resources.SqlResource) (*sql.DB, error) {
	dialect := cfg.SqlDialect()
	connStr := cfg.ConnectionString()

	conn, err := sql.Open(dialect, connStr)
	if err != nil {
		return nil, rerrors.Wrap(err, "error checking connection to postgres")
	}

	closer.Add(func() error {
		return conn.Close()
	})

	goose.SetLogger(sqlLogger{})
	err = goose.SetDialect(dialect)
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting dialect")
	}

	mig := cfg.MigrationFolder()
	if mig == "" {
		mig = "./migrations"
	}

	err = goose.Up(conn, mig)
	if err != nil {
		return nil, rerrors.Wrap(err, "error performing up")
	}

	return conn, nil
}

type DB interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type sqlLogger struct{}

func (s sqlLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (s sqlLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
