package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	infradb "forge-rpg/internal/infrastructure/db"
	"forge-rpg/internal/domain"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := infradb.Open(":memory:")
	if err != nil {
		t.Fatalf("openTestDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func insertNarrativeEntry(t *testing.T, db *sql.DB, id, category, content string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO narrative_entries (id, category, content, created_at) VALUES (?, ?, ?, ?)`,
		id, category, content, time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		t.Fatalf("insertNarrativeEntry: %v", err)
	}
}

func insertNarrativeCompat(t *testing.T, db *sql.DB, entryID, dimension, value, group string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO narrative_compatibility (entry_id, dimension, value, group_name) VALUES (?, ?, ?, ?)`,
		entryID, dimension, value, group,
	)
	if err != nil {
		t.Fatalf("insertNarrativeCompat %q %q %q: %v", entryID, dimension, value, err)
	}
}

// TestNarrativeRepository_FindByCategory_DefaultWeight verifies that entries
// with no compatibility rows get weight 2 (universal).
func TestNarrativeRepository_FindByCategory_DefaultWeight(t *testing.T) {
	db := openTestDB(t)
	insertNarrativeEntry(t, db, "e1", "background", "Contenido universal")
	// No compat rows → default weight 2

	repo := infradb.NewNarrativeRepository(db)
	entries, err := repo.FindByCategory(context.Background(),
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Weight != 2 {
		t.Errorf("expected weight 2, got %d", entries[0].Weight)
	}
	if entries[0].Block.Content != "Contenido universal" {
		t.Errorf("unexpected content: %q", entries[0].Block.Content)
	}
}

// TestNarrativeRepository_FindByCategory_ExcludedNeverReturned verifies that
// entries marked as excluded for the queried class are omitted from results.
func TestNarrativeRepository_FindByCategory_ExcludedNeverReturned(t *testing.T) {
	db := openTestDB(t)
	insertNarrativeEntry(t, db, "e-excluded", "background", "Contenido excluido")
	insertNarrativeCompat(t, db, "e-excluded", "class", "fighter", "excluded")
	insertNarrativeEntry(t, db, "e-allowed", "background", "Contenido permitido")
	// e-allowed has no compat rows → weight 2

	repo := infradb.NewNarrativeRepository(db)
	entries, err := repo.FindByCategory(context.Background(),
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}

	for _, e := range entries {
		if e.Block.Content == "Contenido excluido" {
			t.Error("excluded entry appeared in results")
		}
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry (the allowed one), got %d", len(entries))
	}
}

// TestNarrativeRepository_FindByCategory_PrimaryWeight verifies that entries
// with primary compatibility have weight 10.
func TestNarrativeRepository_FindByCategory_PrimaryWeight(t *testing.T) {
	db := openTestDB(t)
	insertNarrativeEntry(t, db, "e-primary", "background", "Contenido primario")
	insertNarrativeCompat(t, db, "e-primary", "class", "fighter", "primary")

	repo := infradb.NewNarrativeRepository(db)
	entries, err := repo.FindByCategory(context.Background(),
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Weight != 10 {
		t.Errorf("expected weight 10, got %d", entries[0].Weight)
	}
}

// TestNarrativeRepository_FindByCategory_ExcludedWinsOverPrimary verifies that
// if an entry is primary for species but excluded for class, excluded (0) wins.
func TestNarrativeRepository_FindByCategory_ExcludedWinsOverPrimary(t *testing.T) {
	db := openTestDB(t)
	insertNarrativeEntry(t, db, "e-conflict", "background", "Conflicto de peso")
	insertNarrativeCompat(t, db, "e-conflict", "class", "fighter", "excluded")
	insertNarrativeCompat(t, db, "e-conflict", "species", "human", "primary")

	repo := infradb.NewNarrativeRepository(db)
	entries, err := repo.FindByCategory(context.Background(),
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}
	// Excluded should win — entry must not appear
	for _, e := range entries {
		if e.Block.Content == "Conflicto de peso" {
			t.Error("entry should be excluded when class=excluded, but appeared in results")
		}
	}
}

// TestNarrativeRepository_FindByCategory_StableOrdering verifies that two
// calls with identical inputs return entries in the same order.
func TestNarrativeRepository_FindByCategory_StableOrdering(t *testing.T) {
	db := openTestDB(t)
	for i := 0; i < 10; i++ {
		insertNarrativeEntry(t, db,
			"e-"+string(rune('a'+i)),
			"background",
			"Entrada "+string(rune('A'+i)),
		)
	}

	repo := infradb.NewNarrativeRepository(db)
	ctx := context.Background()
	entries1, err := repo.FindByCategory(ctx,
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	entries2, err := repo.FindByCategory(ctx,
		domain.NarrativeBackground, domain.ClassFighter, domain.SpeciesHuman)
	if err != nil {
		t.Fatalf("second call: %v", err)
	}

	if len(entries1) != len(entries2) {
		t.Fatalf("len mismatch: %d vs %d", len(entries1), len(entries2))
	}
	for i := range entries1 {
		if entries1[i].Block.Content != entries2[i].Block.Content {
			t.Errorf("entry %d differs: %q vs %q",
				i, entries1[i].Block.Content, entries2[i].Block.Content)
		}
	}
}

// TestNarrativeRepository_Count verifies Count returns the correct row count.
func TestNarrativeRepository_Count(t *testing.T) {
	db := openTestDB(t)
	insertNarrativeEntry(t, db, "n1", "background", "A")
	insertNarrativeEntry(t, db, "n2", "motivation", "B")
	insertNarrativeEntry(t, db, "n3", "secret", "C")

	repo := infradb.NewNarrativeRepository(db)
	n, err := repo.Count(context.Background())
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}
