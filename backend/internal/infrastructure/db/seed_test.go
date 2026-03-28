package db_test

import (
	"context"
	"testing"

	"forge-rpg/internal/domain"
	infradb "forge-rpg/internal/infrastructure/db"
)

// TestSeedContentIfEmpty_EmptyDB verifies that a fresh DB gets seeded
// to at least 200 narrative entries and names are populated.
func TestSeedContentIfEmpty_EmptyDB(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	// Verify narrative entries
	var narrativeCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries`).Scan(&narrativeCount); err != nil {
		t.Fatalf("count narrative_entries: %v", err)
	}
	if narrativeCount < 200 {
		t.Errorf("expected ≥200 narrative entries, got %d", narrativeCount)
	}

	// Verify name entries exist
	var nameCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&nameCount); err != nil {
		t.Fatalf("count name_entries: %v", err)
	}
	if nameCount == 0 {
		t.Error("name_entries is empty after seed")
	}
}

// TestSeedContentIfEmpty_Idempotent verifies that calling seed twice
// produces no duplicates and the same row count.
func TestSeedContentIfEmpty_Idempotent(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("first SeedContentIfEmpty: %v", err)
	}

	var countAfterFirst int
	db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries`).Scan(&countAfterFirst) //nolint:errcheck

	// Second call — must be a no-op
	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("second SeedContentIfEmpty: %v", err)
	}

	var countAfterSecond int
	db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries`).Scan(&countAfterSecond) //nolint:errcheck

	if countAfterFirst != countAfterSecond {
		t.Errorf("second seed changed row count: %d → %d", countAfterFirst, countAfterSecond)
	}
}

// TestSeedContentIfEmpty_AllCategoriesPresent verifies that all three
// narrative categories are represented in the seed data.
func TestSeedContentIfEmpty_AllCategoriesPresent(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	for _, cat := range []string{"background", "motivation", "secret"} {
		var n int
		db.QueryRowContext(ctx, `SELECT COUNT(*) FROM narrative_entries WHERE category = ?`, cat).Scan(&n) //nolint:errcheck
		if n < 50 {
			t.Errorf("category %q has only %d entries (expected ≥50)", cat, n)
		}
	}
}

// TestSeedNarrativeByVersion_ReachesV3 verifies that narrative_version = 3 after
// SeedContentIfEmpty completes on a fresh DB.
func TestSeedNarrativeByVersion_ReachesV3(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	var version int
	if err := db.QueryRowContext(ctx,
		`SELECT narrative_version FROM seed_version WHERE id = 1`,
	).Scan(&version); err != nil {
		t.Fatalf("read narrative_version: %v", err)
	}
	if version != 3 {
		t.Errorf("expected narrative_version = 3, got %d", version)
	}
}

// TestSeedNarrativeByVersion_GnomeMotivationExcludesHalfOrc is a regression test for
// the coherence bug where a Half-Orc received a gnome-tagged motivation entry.
// Verifies that species-exclusion compat rows are correctly seeded and filtered.
func TestSeedNarrativeByVersion_GnomeMotivationExcludesHalfOrc(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	repo := infradb.NewNarrativeRepository(db)
	entries, err := repo.FindByCategory(ctx,
		domain.NarrativeMotivation, domain.ClassArtificer, domain.SpeciesHalfOrc)
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}

	for _, e := range entries {
		if e.Block.Content == "Tu curiosidad gnoma no tiene límites: necesitás entender cómo funciona todo, desarmarlo si es necesario, y armarlo de nuevo pero mejor." {
			t.Error("gnome-exclusive motivation appeared in half-orc pool — species exclusion not applied")
		}
	}
}

// TestSeedContentIfEmpty_NameCoveragePerSpecies verifies that major species
// keys have ≥50 names per gender after seeding.
func TestSeedContentIfEmpty_NameCoveragePerSpecies(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	keys := []string{
		"human", "high-elf", "wood-elf", "drow",
		"hill-dwarf", "mountain-dwarf",
		"lightfoot", "stout",
		"forest-gnome", "rock-gnome",
		"half-elf", "half-orc", "tiefling", "dragonborn",
	}

	for _, key := range keys {
		for _, gender := range []string{"male", "female"} {
			var n int
			db.QueryRowContext(ctx,
				`SELECT COUNT(*) FROM name_entries WHERE species_key = ? AND gender = ?`,
				key, gender,
			).Scan(&n) //nolint:errcheck
			if n < 50 {
				t.Errorf("species %q gender %q: only %d names (expected ≥50)", key, gender, n)
			}
		}
	}
}
