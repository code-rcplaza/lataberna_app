package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a SQLite database at the given path with foreign keys enabled
// and WAL journal mode for better concurrency. Creates tables if they don't exist.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("db.migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         TEXT PRIMARY KEY,
			email      TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id         TEXT PRIMARY KEY,
			user_id    TEXT NOT NULL REFERENCES users(id),
			expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS magic_link_tokens (
			id            TEXT PRIMARY KEY,
			hashed_token  TEXT NOT NULL UNIQUE,
			email         TEXT NOT NULL,
			expires_at    TEXT NOT NULL,
			used_at       TEXT,
			created_at    TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS characters (
			id                  TEXT PRIMARY KEY,
			user_id             TEXT NOT NULL REFERENCES users(id),
			name                TEXT NOT NULL,
			species             TEXT NOT NULL,
			sub_species         TEXT,
			class               TEXT NOT NULL,
			level               INTEGER NOT NULL DEFAULT 1,
			ruleset             TEXT NOT NULL DEFAULT '5e',
			ability_bonus_source TEXT NOT NULL DEFAULT 'species',
			base_stats          TEXT NOT NULL,
			final_stats         TEXT NOT NULL,
			modifiers           TEXT NOT NULL,
			derived             TEXT NOT NULL,
			background          TEXT NOT NULL,
			motivation          TEXT NOT NULL,
			secret              TEXT NOT NULL,
			locks               TEXT NOT NULL,
			seed                INTEGER,
			created_at          TEXT NOT NULL,
			updated_at          TEXT NOT NULL
		);
	`)
	return err
}
