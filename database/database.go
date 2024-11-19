package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// import "postgres" driver without directly referencing the library
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var Db *sql.DB

const maxOpenDBConns = 10
const maxIdleDBConns = 5
const connMaxDBLifetime = 5 * time.Minute

func InitDB(dbString string) error {
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxOpenDBConns)
	db.SetMaxIdleConns(maxIdleDBConns)
	db.SetConnMaxLifetime(connMaxDBLifetime)

	// test db is connected
	if err = db.Ping(); err != nil {
		return err
	}

	if err = MigrateUp(); err != nil {
		return err
	}

	Db = db
	return nil
}

func CloseDB() error {
	return Db.Close()
}

// auto migrate on server startup
func MigrateUp() error {
	m, err := migrate.New(
		"file://database/migrations",
		os.Getenv("POSTGRESQL_URL"),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get current migration version
	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to retrieve migration version: %w", err)
	}

	// Log migration status
	if err == migrate.ErrNoChange {
		fmt.Println("No new migrations applied")
	} else {
		fmt.Printf("Current Migration version: %d\n", version)
	}
	if dirty {
		fmt.Println("Warning: The database is in a dirty state!")
	}

	return nil
}
