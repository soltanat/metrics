package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Conn interface {
	Ping(ctx context.Context) error
}

func New(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}
