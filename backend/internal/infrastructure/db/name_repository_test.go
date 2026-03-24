package db_test

import (
	"context"
	"testing"
	"time"

	infradb "forge-rpg/internal/infrastructure/db"
)

func insertNameEntry(t *testing.T, db interface {
	Exec(string, ...any) (interface{ LastInsertId() (int64, error); RowsAffected() (int64, error) }, error)
}, id, speciesKey, gender, name string) {
	// Use a helper that avoids importing sql directly in the test
	t.Helper()
}

// TestNameRepository_FindBySpeciesGender verifies correct filtering.
func TestNameRepository_FindBySpeciesGender(t *testing.T) {
	db := openTestDB(t)

	now := time.Now().UTC().Format(time.RFC3339Nano)
	for i, name := range []string{"Aldric", "Brennan", "Cael"} {
		id := "nm-" + string(rune('a'+i))
		db.Exec(`INSERT INTO name_entries (id, species_key, gender, name, created_at) VALUES (?,?,?,?,?)`,
			id, "human", "male", name, now)
	}
	for i, name := range []string{"Aelara", "Brenna"} {
		id := "nf-" + string(rune('a'+i))
		db.Exec(`INSERT INTO name_entries (id, species_key, gender, name, created_at) VALUES (?,?,?,?,?)`,
			id, "human", "female", name, now)
	}
	// Add elf names — must not appear in human queries
	db.Exec(`INSERT INTO name_entries (id, species_key, gender, name, created_at) VALUES (?,?,?,?,?)`,
		"ne-a", "high-elf", "male", "Aelindor", now)

	repo := infradb.NewNameRepository(db)

	t.Run("male human names", func(t *testing.T) {
		names, err := repo.FindBySpeciesGender(context.Background(), "human", "male")
		if err != nil {
			t.Fatalf("FindBySpeciesGender: %v", err)
		}
		if len(names) != 3 {
			t.Errorf("expected 3 male human names, got %d: %v", len(names), names)
		}
	})

	t.Run("female human names", func(t *testing.T) {
		names, err := repo.FindBySpeciesGender(context.Background(), "human", "female")
		if err != nil {
			t.Fatalf("FindBySpeciesGender: %v", err)
		}
		if len(names) != 2 {
			t.Errorf("expected 2 female human names, got %d: %v", len(names), names)
		}
	})

	t.Run("unknown species returns empty", func(t *testing.T) {
		names, err := repo.FindBySpeciesGender(context.Background(), "unknown-species", "male")
		if err != nil {
			t.Fatalf("FindBySpeciesGender: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("expected 0 names for unknown species, got %d", len(names))
		}
	})
}

// TestNameRepository_Count verifies Count returns the correct row count.
func TestNameRepository_Count(t *testing.T) {
	db := openTestDB(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	for i, entry := range []struct{ id, key, gender, name string }{
		{"a", "human", "male", "Aldric"},
		{"b", "human", "female", "Aelara"},
		{"c", "high-elf", "male", "Aelindor"},
	} {
		_ = i
		db.Exec(`INSERT INTO name_entries (id, species_key, gender, name, created_at) VALUES (?,?,?,?,?)`,
			entry.id, entry.key, entry.gender, entry.name, now)
	}

	repo := infradb.NewNameRepository(db)
	n, err := repo.Count(context.Background())
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}
