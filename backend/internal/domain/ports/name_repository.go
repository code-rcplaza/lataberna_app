package ports

import "context"

// NameRepository defines the persistence interface for character name pools.
// The infrastructure layer implements this; the domain never imports infrastructure.
type NameRepository interface {
	// FindBySpeciesGender returns all name strings for the given species key and
	// gender, ordered by entry ID for deterministic selection.
	FindBySpeciesGender(
		ctx        context.Context,
		speciesKey string,
		gender     string,
	) ([]string, error)

	// Count returns the total number of name entries in the store.
	// Used for idempotency check in seed loading.
	Count(ctx context.Context) (int, error)
}
