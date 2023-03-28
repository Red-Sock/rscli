package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"financial-microservice/internal/utils/closer"
)

func New(ctx context.Context, connectionString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "error checking connection to redis")
	}

	closer.Add(func() error {
		return conn.Close(ctx)
	})

	return conn, nil
}
