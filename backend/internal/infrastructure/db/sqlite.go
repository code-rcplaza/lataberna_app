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
	// Step 1: Base schema — all tables with IF NOT EXISTS (idempotent).
	// name_entries is created here with legacy schema; v2 migration below upgrades it.
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

		CREATE TABLE IF NOT EXISTS narrative_entries (
			id         TEXT PRIMARY KEY,
			category   TEXT NOT NULL CHECK(category IN ('background','motivation','secret')),
			content    TEXT NOT NULL,
			created_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS narrative_compatibility (
			entry_id   TEXT NOT NULL REFERENCES narrative_entries(id) ON DELETE CASCADE,
			dimension  TEXT NOT NULL CHECK(dimension IN ('class','species')),
			value      TEXT NOT NULL,
			group_name TEXT NOT NULL CHECK(group_name IN ('primary','secondary','excluded')),
			PRIMARY KEY (entry_id, dimension, value)
		);

		CREATE INDEX IF NOT EXISTS idx_narrative_compat
			ON narrative_compatibility(dimension, value, group_name);

		CREATE TABLE IF NOT EXISTS name_entries (
			id          TEXT PRIMARY KEY,
			species_key TEXT NOT NULL,
			gender      TEXT NOT NULL CHECK(gender IN ('male','female')),
			name        TEXT NOT NULL,
			created_at  TEXT NOT NULL,
			UNIQUE(species_key, gender, name)
		);

		CREATE INDEX IF NOT EXISTS idx_name_entries
			ON name_entries(species_key, gender);
	`)
	if err != nil {
		return err
	}

	// Step 2: seed_version table — tracks seed data completion state.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS seed_version (
			id      INTEGER PRIMARY KEY CHECK(id = 1),
			version INTEGER NOT NULL DEFAULT 0
		);
		INSERT OR IGNORE INTO seed_version(id, version) VALUES(1, 0);
	`); err != nil {
		return fmt.Errorf("migrate: seed_version: %w", err)
	}

	// Step 3: Upgrade name_entries to v2 schema (adds name_type column) if needed.
	needs, err := nameEntriesNeedsV2(db)
	if err != nil {
		return fmt.Errorf("migrate: check v2: %w", err)
	}
	if needs {
		if err := migrateNameEntriesV2(db); err != nil {
			return fmt.Errorf("migrate: name_entries v2: %w", err)
		}
	}

	return nil
}

// nameEntriesNeedsV2 reports whether name_entries is missing the name_type column.
func nameEntriesNeedsV2(db *sql.DB) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(name_entries)`)
	if err != nil {
		return false, fmt.Errorf("nameEntriesNeedsV2: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notNull, pk int
		var colName, colType string
		var dfltValue sql.NullString
		if err := rows.Scan(&cid, &colName, &colType, &notNull, &dfltValue, &pk); err != nil {
			return false, fmt.Errorf("nameEntriesNeedsV2: scan: %w", err)
		}
		if colName == "name_type" {
			return false, nil
		}
	}
	return true, rows.Err()
}

// migrateNameEntriesV2 adds the name_type discriminator column to name_entries via
// rename → recreate → copy → drop (SQLite cannot ALTER COLUMN directly).
//
// If the table had existing rows, seed_version is advanced to 1 — those rows are
// preserved with name_type='first_name', so the seeder can skip the v1 phase.
func migrateNameEntriesV2(db *sql.DB) error {
	// Count existing rows before migration to decide seed_version advancement.
	var rowCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM name_entries`).Scan(&rowCount); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: count: %w", err)
	}

	if _, err := db.Exec(`ALTER TABLE name_entries RENAME TO name_entries_v1`); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: rename: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE name_entries (
			id          TEXT PRIMARY KEY,
			species_key TEXT NOT NULL,
			gender      TEXT NOT NULL CHECK(gender IN ('male','female','any')),
			name_type   TEXT NOT NULL DEFAULT 'first_name'
				CHECK(name_type IN ('first_name','surname','clan_name','family_name','nickname','virtue_word','infernal_name')),
			name        TEXT NOT NULL,
			created_at  TEXT NOT NULL,
			UNIQUE(species_key, gender, name_type, name)
		)
	`); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: create: %w", err)
	}

	if _, err := db.Exec(`
		INSERT INTO name_entries
		SELECT id, species_key, gender, 'first_name', name, created_at
		FROM name_entries_v1
	`); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: copy: %w", err)
	}

	if _, err := db.Exec(`DROP TABLE name_entries_v1`); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: drop: %w", err)
	}

	if _, err := db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_name_entries ON name_entries(species_key, gender, name_type)`,
	); err != nil {
		return fmt.Errorf("migrateNameEntriesV2: index: %w", err)
	}

	// Advance seed_version to 1 only when rows existed — they are now typed
	// as 'first_name', so the seeder should skip phase 1 and only run phase 2.
	if rowCount > 0 {
		if _, err := db.Exec(`UPDATE seed_version SET version = 1 WHERE id = 1`); err != nil {
			return fmt.Errorf("migrateNameEntriesV2: seed_version: %w", err)
		}
	}

	return nil
}
