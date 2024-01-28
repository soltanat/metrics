package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(connString string) error {
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) && err != nil {
		log.Fatal(err)
	}

	return nil
}
