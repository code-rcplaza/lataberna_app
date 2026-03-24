package db

import (
	"context"
	"database/sql"
	"errors"
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

// FindByType returns all names matching (species_key, gender, name_type), ordered by ID.
// Returns an error wrapping ports.ErrEmptyNamePool when the result set is empty.
func (r *nameRepository) FindByType(
	ctx        context.Context,
	speciesKey string,
	gender     string,
	nameType   string,
) ([]string, error) {
	const q = `SELECT name FROM name_entries WHERE species_key = ? AND gender = ? AND name_type = ? ORDER BY id`

	rows, err := r.db.QueryContext(ctx, q, speciesKey, gender, nameType)
	if err != nil {
		return nil, fmt.Errorf("nameRepository.FindByType: query: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("nameRepository.FindByType: scan: %w", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("nameRepository.FindByType: rows: %w", err)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf(
			"nameRepository.FindByType: species=%q gender=%q type=%q: %w",
			speciesKey, gender, nameType, ports.ErrEmptyNamePool,
		)
	}
	return names, nil
}

// FindBySpeciesGender returns all first_name entries for the given species key and gender.
// Returns nil, nil when no names are found (backward-compatible with callers checking len).
func (r *nameRepository) FindBySpeciesGender(
	ctx        context.Context,
	speciesKey string,
	gender     string,
) ([]string, error) {
	names, err := r.FindByType(ctx, speciesKey, gender, "first_name")
	if errors.Is(err, ports.ErrEmptyNamePool) {
		return nil, nil
	}
	return names, err
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
