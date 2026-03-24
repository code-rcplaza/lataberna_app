package ports

import (
	"context"

	"forge-rpg/internal/domain"
)

// WeightedNarrativeEntry is a narrative block with an associated weight for
// weighted-random selection. Weight values: primary=10, secondary=4, default=2, excluded=0.
// Entries with weight 0 are filtered out before selection.
type WeightedNarrativeEntry struct {
	Block  domain.NarrativeBlock
	Weight int
}

// NarrativeRepository defines the persistence interface for narrative content.
// The infrastructure layer implements this; the domain never imports infrastructure.
type NarrativeRepository interface {
	// FindByCategory returns all narrative entries for the given category, with
	// weights resolved based on class and species compatibility rows.
	// Entries with weight 0 (excluded) are omitted from the result.
	// Results are ordered by entry ID for deterministic pool construction.
	FindByCategory(
		ctx      context.Context,
		category domain.NarrativeCategory,
		class    domain.Class,
		species  domain.Species,
	) ([]WeightedNarrativeEntry, error)

	// Count returns the total number of narrative entries in the store.
	// Used for idempotency check in seed loading.
	Count(ctx context.Context) (int, error)
}
