package db_test

import (
	"context"
	"testing"

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
