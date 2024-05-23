package database

import (
	"database/sql"
	"fmt"
	"time"

	// import "postgres" driver without directly referencing the library
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

var dbConn = &Database{}

const maxOpenDBConns = 10
const maxIdleDBConns = 5
const connMaxDBLifetime = 5 * time.Minute

func ConnectPostgres(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDBConns)
	db.SetMaxIdleConns(maxIdleDBConns)
	db.SetConnMaxLifetime(connMaxDBLifetime)

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	dbConn.DB = db
	return dbConn, nil
}

func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	fmt.Println("*** Pinged database successfully! ***")
	return nil
}
