package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a SQLite database at the given path with foreign keys enabled
// and WAL journal mode for better concurrency.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
