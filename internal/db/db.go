package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	*pgx.Conn
}

func New(ctx context.Context, connString string) (*DB, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return &DB{conn}, nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.Close(ctx)
}

func (db *DB) Ping(ctx context.Context) error {
	return db.Ping(ctx)
}
