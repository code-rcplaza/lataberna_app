package db

import (
	"context"
	"database/sql"
	"fmt"

	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.NameRepository = (*nameRepository)(nil)

// nameRepository implements ports.NameRepository using SQLite.
type nameRepository struct {
	db *sql.DB
}

// NewNameRepository creates a new nameRepository.
func NewNameRepository(db *sql.DB) *nameRepository {
	return &nameRepository{db: db}
}

// FindBySpeciesGender returns all names for the given species key and gender,
// ordered by ID for deterministic selection with a seeded RNG.
func (r *nameRepository) FindBySpeciesGender(
	ctx        context.Context,
	speciesKey string,
	gender     string,
) ([]string, error) {
	const q = `SELECT name FROM name_entries WHERE species_key = ? AND gender = ? ORDER BY id`

	rows, err := r.db.QueryContext(ctx, q, speciesKey, gender)
	if err != nil {
		return nil, fmt.Errorf("nameRepository.FindBySpeciesGender: query: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("nameRepository.FindBySpeciesGender: scan: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("nameRepository.FindBySpeciesGender: rows: %w", err)
	}
	return names, nil
}

// Count returns the total number of name entries.
func (r *nameRepository) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("nameRepository.Count: %w", err)
	}
	return n, nil
}
