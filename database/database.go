package database

import (
	"database/sql"
	"time"

	// import "postgres" driver without directly referencing the library
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

	Db = db
	return nil
}

func CloseDB() error {
	return Db.Close()
}

// TODO: set up auto migrate on server startup
func Migrate() {

}
