package ports

import (
	"context"
	"errors"
)

// ErrEmptyNamePool is returned by NameRepository.FindByType when the queried pool
// contains no entries for the given (speciesKey, gender, nameType) triple.
var ErrEmptyNamePool = errors.New("empty name pool")

// NameRepository defines the persistence interface for character name pools.
// The infrastructure layer implements this; the domain never imports infrastructure.
type NameRepository interface {
	// FindBySpeciesGender returns all first_name strings for the given species key
	// and gender, ordered by entry ID for deterministic selection.
	// Returns nil, nil when no names are found (backward-compatible).
	FindBySpeciesGender(
		ctx        context.Context,
		speciesKey string,
		gender     string,
	) ([]string, error)

	// FindByType returns all name strings for the given species key, gender, and
	// name_type, ordered by entry ID. Returns an error wrapping ErrEmptyNamePool
	// when the result set is empty.
	FindByType(
		ctx        context.Context,
		speciesKey string,
		gender     string,
		nameType   string,
	) ([]string, error)

	// Count returns the total number of name entries in the store.
	// Used for idempotency check in seed loading.
	Count(ctx context.Context) (int, error)
}
