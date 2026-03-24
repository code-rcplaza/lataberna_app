package db

import (
	"context"
	"database/sql"
	"fmt"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.NarrativeRepository = (*narrativeRepository)(nil)

// narrativeRepository implements ports.NarrativeRepository using SQLite.
//
// Weight resolution via SQL:
//   - LEFT JOIN narrative_compatibility to get all compat rows for this entry
//     where dimension is 'class' or 'species' and value matches the requested class/species
//   - MIN(weight) across matching rows: excluded (0) always wins over primary (10)
//   - Entries with no matching compat rows get weight 2 (universal default)
//   - HAVING weight > 0 filters out excluded entries from the pool
//   - ORDER BY e.id guarantees stable pool ordering for seed reproducibility
type narrativeRepository struct {
	db *sql.DB
}

// NewNarrativeRepository creates a new narrativeRepository.
func NewNarrativeRepository(db *sql.DB) *narrativeRepository {
	return &narrativeRepository{db: db}
}

// FindByCategory returns weighted narrative entries for the given category, class, and species.
// Entries excluded for the queried class or species are omitted (weight = 0 never returned).
func (r *narrativeRepository) FindByCategory(
	ctx      context.Context,
	category domain.NarrativeCategory,
	class    domain.Class,
	species  domain.Species,
) ([]ports.WeightedNarrativeEntry, error) {
	const q = `
SELECT e.id, e.content,
    COALESCE(MIN(CASE c.group_name
        WHEN 'excluded'   THEN 0
        WHEN 'primary'    THEN 10
        WHEN 'secondary'  THEN 4
        ELSE 2 END), 2) AS weight
FROM narrative_entries e
LEFT JOIN narrative_compatibility c
    ON c.entry_id = e.id
    AND c.dimension IN ('class', 'species')
    AND c.value IN (?, ?)
WHERE e.category = ?
GROUP BY e.id, e.content
HAVING weight > 0
ORDER BY e.id`

	rows, err := r.db.QueryContext(ctx, q, string(class), string(species), string(category))
	if err != nil {
		return nil, fmt.Errorf("narrativeRepository.FindByCategory: query: %w", err)
	}
	defer rows.Close()

	var out []ports.WeightedNarrativeEntry
	for rows.Next() {
		var id, content string
		var weight int
		if err := rows.Scan(&id, &content, &weight); err != nil {
			return nil, fmt.Errorf("narrativeRepository.FindByCategory: scan: %w", err)
		}
		out = append(out, ports.WeightedNarrativeEntry{
			Block: domain.NarrativeBlock{
				Category: category,
				Content:  content,
				Tags:     nil, // tags live in compatibility rows, not in the block
			},
			Weight: weight,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("narrativeRepository.FindByCategory: rows: %w", err)
	}
	return out, nil
}

// Count returns the total number of narrative entries.
func (r *narrativeRepository) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries`).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("narrativeRepository.Count: %w", err)
	}
	return n, nil
}
