package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"financial-microservice/internal/config"
	"financial-microservice/internal/utils/closer"
)

func New(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, createConnectionString(cfg))
	if err != nil {
		return nil, errors.Wrap(err, "error checking connection to redis")
	}

	closer.Add(func() error {
		return conn.Close(ctx)
	})

	return conn, nil
}

func createConnectionString(cfg *config.Config) string {
	sslMode := cfg.GetString(config.DataSourcesPostgresSslmode)

	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.GetString(config.DataSourcesPostgresUser),
		cfg.GetString(config.DataSourcesPostgresPwd),
		cfg.GetString(config.DataSourcesPostgresHost),
		cfg.GetString(config.DataSourcesPostgresPort),
		cfg.GetString(config.DataSourcesPostgresName),
		sslMode,
	)
}
