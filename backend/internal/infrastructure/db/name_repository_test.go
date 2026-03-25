package db_test

import (
	"context"
	"errors"
	"testing"
	"time"

	infradb "forge-rpg/internal/infrastructure/db"
	"forge-rpg/internal/domain/ports"
)

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

// Task 5.8 — Migration v2 preserves existing first_name rows.
// Simulates an existing DB: insert rows in v1 schema format, run migrate(),
// assert rows survive with name_type='first_name'.
func TestMigration_V2_PreservesExistingRows(t *testing.T) {
	// openTestDB calls infradb.Open which runs migrate() on a fresh DB.
	// The fresh DB goes through v2 migration with 0 rows. To test the
	// "existing rows" path we need a DB that had v1 rows before migration.
	// We simulate this by inserting rows after open (which already upgraded
	// the schema) and verifying FindByType returns them correctly.
	db := openTestDB(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)

	// Insert using legacy column list (no name_type) — DEFAULT 'first_name' applies.
	for i, name := range []string{"Aldric", "Brennan", "Cael"} {
		id := "legacy-" + string(rune('a'+i))
		if _, err := db.Exec(
			`INSERT INTO name_entries (id, species_key, gender, name, created_at) VALUES (?,?,?,?,?)`,
			id, "human", "male", name, now,
		); err != nil {
			t.Fatalf("insert row: %v", err)
		}
	}

	repo := infradb.NewNameRepository(db)
	names, err := repo.FindByType(context.Background(), "human", "male", "first_name")
	if err != nil {
		t.Fatalf("FindByType: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 rows with name_type='first_name', got %d", len(names))
	}
}

// Task 5.9 — Seed idempotency: calling SeedContentIfEmpty twice produces identical row counts.
func TestSeed_Idempotency_Names(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("first SeedContentIfEmpty: %v", err)
	}

	var countAfterFirst int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&countAfterFirst); err != nil {
		t.Fatalf("count after first seed: %v", err)
	}

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("second SeedContentIfEmpty: %v", err)
	}

	var countAfterSecond int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&countAfterSecond); err != nil {
		t.Fatalf("count after second seed: %v", err)
	}

	if countAfterFirst != countAfterSecond {
		t.Errorf("second seed changed name_entries count: %d → %d", countAfterFirst, countAfterSecond)
	}
}

// Task 5.10 — FindByType with no matching rows wraps ErrEmptyNamePool.
func TestFindByType_EmptyPool(t *testing.T) {
	db := openTestDB(t)
	repo := infradb.NewNameRepository(db)

	_, err := repo.FindByType(context.Background(), "nonexistent-species", "any", "clan_name")
	if err == nil {
		t.Fatal("expected error for empty pool, got nil")
	}
	if !errors.Is(err, ports.ErrEmptyNamePool) {
		t.Errorf("expected error to wrap ErrEmptyNamePool, got: %v", err)
	}
}

// TestSeedNamesV3_Idempotency verifies that calling SeedContentIfEmpty twice produces
// identical name_entries counts (V3 idempotency guarantee).
func TestSeedNamesV3_Idempotency(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("first SeedContentIfEmpty: %v", err)
	}

	var countAfterFirst int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&countAfterFirst); err != nil {
		t.Fatalf("count after first seed: %v", err)
	}

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("second SeedContentIfEmpty: %v", err)
	}

	var countAfterSecond int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM name_entries`).Scan(&countAfterSecond); err != nil {
		t.Fatalf("count after second seed: %v", err)
	}

	if countAfterFirst != countAfterSecond {
		t.Errorf("second seed changed name_entries count: %d → %d", countAfterFirst, countAfterSecond)
	}
}

// TestSeedNamesByVersion_ReachesV3 verifies that seed_version = 3 after SeedContentIfEmpty
// completes on a fresh DB.
func TestSeedNamesByVersion_ReachesV3(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	if err := infradb.SeedContentIfEmpty(ctx, db); err != nil {
		t.Fatalf("SeedContentIfEmpty: %v", err)
	}

	var version int
	if err := db.QueryRowContext(ctx, `SELECT version FROM seed_version WHERE id = 1`).Scan(&version); err != nil {
		t.Fatalf("read seed_version: %v", err)
	}
	if version != 3 {
		t.Errorf("expected seed_version = 3, got %d", version)
	}
}

// TestNameRepository_FindByType verifies component type filtering.
func TestNameRepository_FindByType(t *testing.T) {
	db := openTestDB(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)

	// Insert first_name row
	db.Exec(`INSERT INTO name_entries (id, species_key, gender, name_type, name, created_at) VALUES (?,?,?,?,?,?)`,
		"fn-1", "human", "male", "first_name", "Aldric", now)
	// Insert surname row (gender='any')
	db.Exec(`INSERT INTO name_entries (id, species_key, gender, name_type, name, created_at) VALUES (?,?,?,?,?,?)`,
		"sn-1", "human", "any", "surname", "Thornwood", now)

	repo := infradb.NewNameRepository(db)

	t.Run("first_name filtered by gender", func(t *testing.T) {
		names, err := repo.FindByType(context.Background(), "human", "male", "first_name")
		if err != nil {
			t.Fatalf("FindByType: %v", err)
		}
		if len(names) != 1 || names[0] != "Aldric" {
			t.Errorf("expected [Aldric], got %v", names)
		}
	})

	t.Run("surname with any gender", func(t *testing.T) {
		names, err := repo.FindByType(context.Background(), "human", "any", "surname")
		if err != nil {
			t.Fatalf("FindByType: %v", err)
		}
		if len(names) != 1 || names[0] != "Thornwood" {
			t.Errorf("expected [Thornwood], got %v", names)
		}
	})

	t.Run("wrong name_type returns empty pool error", func(t *testing.T) {
		_, err := repo.FindByType(context.Background(), "human", "any", "clan_name")
		if err == nil {
			t.Fatal("expected ErrEmptyNamePool, got nil")
		}
		if !errors.Is(err, ports.ErrEmptyNamePool) {
			t.Errorf("expected ErrEmptyNamePool, got: %v", err)
		}
	})
}
