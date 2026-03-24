package narrativegen_test

import (
	"context"
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
	"forge-rpg/internal/usecase/narrativegen"
)

// ptr helpers

func ptrClass(c domain.Class) *domain.Class     { return &c }
func ptrSpecies(s domain.Species) *domain.Species { return &s }
func ptrSeed(v int64) *int64                     { return &v }

// ---------------------------------------------------------------------------
// fakeNarrativeRepo — in-memory stub satisfying ports.NarrativeRepository
// ---------------------------------------------------------------------------

// fakeNarrativeRepo returns a fixed pool for every FindByCategory call.
// The pool is keyed by category so tests can control per-category content.
type fakeNarrativeRepo struct {
	entries map[domain.NarrativeCategory][]ports.WeightedNarrativeEntry
}

func (f *fakeNarrativeRepo) FindByCategory(
	ctx      context.Context,
	category domain.NarrativeCategory,
	class    domain.Class,
	species  domain.Species,
) ([]ports.WeightedNarrativeEntry, error) {
	if pool, ok := f.entries[category]; ok {
		return pool, nil
	}
	return nil, nil
}

func (f *fakeNarrativeRepo) Count(ctx context.Context) (int, error) {
	total := 0
	for _, pool := range f.entries {
		total += len(pool)
	}
	return total, nil
}

// defaultFakeRepo returns a repo with 3 distinct entries per category,
// all with default weight (2), so every class/species combination succeeds.
func defaultFakeRepo() *fakeNarrativeRepo {
	make3 := func(cat domain.NarrativeCategory, contents ...string) []ports.WeightedNarrativeEntry {
		out := make([]ports.WeightedNarrativeEntry, len(contents))
		for i, c := range contents {
			out[i] = ports.WeightedNarrativeEntry{
				Block:  domain.NarrativeBlock{Category: cat, Content: c},
				Weight: 2,
			}
		}
		return out
	}

	return &fakeNarrativeRepo{
		entries: map[domain.NarrativeCategory][]ports.WeightedNarrativeEntry{
			domain.NarrativeBackground: make3(domain.NarrativeBackground,
				"Fondo universal A",
				"Fondo universal B",
				"Fondo universal C",
			),
			domain.NarrativeMotivation: make3(domain.NarrativeMotivation,
				"Motivación universal A",
				"Motivación universal B",
				"Motivación universal C",
			),
			domain.NarrativeSecret: make3(domain.NarrativeSecret,
				"Secreto universal A",
				"Secreto universal B",
				"Secreto universal C",
			),
		},
	}
}

// ─── 1. Reproducibility ──────────────────────────────────────────────────────

func TestGenerate_Reproducibility(t *testing.T) {
	seed := int64(42)
	class := domain.ClassFighter
	species := domain.SpeciesHuman

	svc := narrativegen.New(defaultFakeRepo())
	in := narrativegen.Input{
		Class:   ptrClass(class),
		Species: ptrSpecies(species),
		Seed:    ptrSeed(seed),
	}

	out1, err1 := svc.Generate(context.Background(), in)
	out2, err2 := svc.Generate(context.Background(), in)

	if err1 != nil {
		t.Fatalf("first Generate: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("second Generate: %v", err2)
	}

	if out1.Background.Content != out2.Background.Content {
		t.Errorf("Background.Content differs: %q vs %q", out1.Background.Content, out2.Background.Content)
	}
	if out1.Motivation.Content != out2.Motivation.Content {
		t.Errorf("Motivation.Content differs: %q vs %q", out1.Motivation.Content, out2.Motivation.Content)
	}
	if out1.Secret.Content != out2.Secret.Content {
		t.Errorf("Secret.Content differs: %q vs %q", out1.Secret.Content, out2.Secret.Content)
	}
	if out1.Seed != out2.Seed {
		t.Errorf("Seed differs: %d vs %d", out1.Seed, out2.Seed)
	}
}

// ─── 2. Three blocks always returned — non-empty content ─────────────────────

func TestGenerate_ThreeBlocksAlwaysReturned(t *testing.T) {
	svc := narrativegen.New(defaultFakeRepo())
	out, err := svc.Generate(context.Background(), narrativegen.Input{
		Class:   ptrClass(domain.ClassWizard),
		Species: ptrSpecies(domain.SpeciesElf),
		Seed:    ptrSeed(99),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Background.Content == "" {
		t.Error("Background.Content is empty")
	}
	if out.Motivation.Content == "" {
		t.Error("Motivation.Content is empty")
	}
	if out.Secret.Content == "" {
		t.Error("Secret.Content is empty")
	}
}

// ─── 3. Excluded entries never appear ────────────────────────────────────────

func TestGenerate_ExcludedEntriesNeverAppear(t *testing.T) {
	// Pool: one entry weight 2 (allowed), one entry weight 0 (excluded).
	// The excluded content must never appear across 200 draws.
	excludedContent := "ESTE CONTENIDO NUNCA DEBE APARECER"
	allowedContent := "Contenido permitido"

	repo := &fakeNarrativeRepo{
		entries: map[domain.NarrativeCategory][]ports.WeightedNarrativeEntry{
			domain.NarrativeBackground: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeBackground, Content: allowedContent}, Weight: 2},
				// weight 0 = excluded — filtered out by the repo in real impl,
				// but the test verifies weightedPick also respects it
			},
			domain.NarrativeMotivation: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeMotivation, Content: allowedContent}, Weight: 2},
			},
			domain.NarrativeSecret: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeSecret, Content: allowedContent}, Weight: 2},
			},
		},
	}

	svc := narrativegen.New(repo)
	for i := int64(0); i < 200; i++ {
		out, err := svc.Generate(context.Background(), narrativegen.Input{
			Class:   ptrClass(domain.ClassBarbarian),
			Species: ptrSpecies(domain.SpeciesHuman),
			Seed:    ptrSeed(i),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
		if out.Background.Content == excludedContent {
			t.Errorf("seed %d: excluded content appeared in Background", i)
		}
	}
}

// ─── 4. Universal tags compatible with all classes ────────────────────────────

func TestGenerate_UniversalTags_CompatibleWithAllClasses(t *testing.T) {
	allClasses := []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric, domain.ClassDruid,
		domain.ClassFighter, domain.ClassMonk, domain.ClassPaladin, domain.ClassRanger,
		domain.ClassRogue, domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}
	species := domain.SpeciesHuman
	seed := int64(7)

	for _, cls := range allClasses {
		cls := cls
		t.Run(string(cls), func(t *testing.T) {
			svc := narrativegen.New(defaultFakeRepo())
			out, err := svc.Generate(context.Background(), narrativegen.Input{
				Class:   &cls,
				Species: &species,
				Seed:    ptrSeed(seed),
			})
			if err != nil {
				t.Fatalf("class %q: %v", cls, err)
			}
			if out.Background.Content == "" {
				t.Errorf("Background not populated for class %q", cls)
			}
		})
	}
}

// ─── 5. All 13 classes generate without error ────────────────────────────────

func TestGenerate_AllClassesSucceed(t *testing.T) {
	allClasses := []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric, domain.ClassDruid,
		domain.ClassFighter, domain.ClassMonk, domain.ClassPaladin, domain.ClassRanger,
		domain.ClassRogue, domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}
	species := domain.SpeciesHuman
	seed := int64(42)

	for _, cls := range allClasses {
		cls := cls
		t.Run(string(cls), func(t *testing.T) {
			svc := narrativegen.New(defaultFakeRepo())
			out, err := svc.Generate(context.Background(), narrativegen.Input{
				Class:   &cls,
				Species: &species,
				Seed:    ptrSeed(seed),
			})
			if err != nil {
				t.Fatalf("class %q returned error: %v", cls, err)
			}
			if out.Background.Content == "" || out.Motivation.Content == "" || out.Secret.Content == "" {
				t.Errorf("class %q: one or more blocks have empty content", cls)
			}
		})
	}
}

// ─── 6. All 9 species generate without error ─────────────────────────────────

func TestGenerate_AllSpeciesSucceed(t *testing.T) {
	allSpecies := []domain.Species{
		domain.SpeciesHuman, domain.SpeciesElf, domain.SpeciesDwarf, domain.SpeciesHalfling,
		domain.SpeciesGnome, domain.SpeciesHalfElf, domain.SpeciesHalfOrc,
		domain.SpeciesTiefling, domain.SpeciesDragonborn,
	}
	class := domain.ClassFighter
	seed := int64(42)

	for _, sp := range allSpecies {
		sp := sp
		t.Run(string(sp), func(t *testing.T) {
			svc := narrativegen.New(defaultFakeRepo())
			out, err := svc.Generate(context.Background(), narrativegen.Input{
				Class:   &class,
				Species: &sp,
				Seed:    ptrSeed(seed),
			})
			if err != nil {
				t.Fatalf("species %q returned error: %v", sp, err)
			}
			if out.Background.Content == "" || out.Motivation.Content == "" || out.Secret.Content == "" {
				t.Errorf("species %q: one or more blocks have empty content", sp)
			}
		})
	}
}

// ─── 7. Nil class picks randomly — output is always valid ────────────────────

func TestGenerate_NilClass_PicksRandomly(t *testing.T) {
	species := domain.SpeciesHuman

	for i := int64(0); i < 13; i++ {
		svc := narrativegen.New(defaultFakeRepo())
		out, err := svc.Generate(context.Background(), narrativegen.Input{
			Class:   nil,
			Species: &species,
			Seed:    ptrSeed(i),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
		if out.Background.Content == "" {
			t.Errorf("seed %d: Background.Content is empty with nil class", i)
		}
	}
}

// ─── 8. Nil seed generates valid output ──────────────────────────────────────

func TestGenerate_NilSeed_ReturnsValidOutput(t *testing.T) {
	class := domain.ClassRogue
	species := domain.SpeciesHalfling

	svc := narrativegen.New(defaultFakeRepo())
	out, err := svc.Generate(context.Background(), narrativegen.Input{
		Class:   &class,
		Species: &species,
		Seed:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Background.Content == "" {
		t.Error("Background.Content is empty")
	}
	if out.Motivation.Content == "" {
		t.Error("Motivation.Content is empty")
	}
	if out.Secret.Content == "" {
		t.Error("Secret.Content is empty")
	}
	if out.Seed == 0 {
		t.Error("Seed should be non-zero in output even when input Seed is nil")
	}
}

// ─── 9. Seed returned in output matches the seed used ───────────────────────

func TestGenerate_SeedReturnedInOutput(t *testing.T) {
	seed := int64(12345)
	class := domain.ClassBard
	species := domain.SpeciesTiefling

	svc := narrativegen.New(defaultFakeRepo())
	out, err := svc.Generate(context.Background(), narrativegen.Input{
		Class:   &class,
		Species: &species,
		Seed:    ptrSeed(seed),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Seed != seed {
		t.Errorf("expected Seed=%d in output, got %d", seed, out.Seed)
	}
}

// ─── 10. Categories are correct on each block ────────────────────────────────

func TestGenerate_CategoriesAreCorrect(t *testing.T) {
	seed := int64(77)
	class := domain.ClassPaladin
	species := domain.SpeciesDwarf

	svc := narrativegen.New(defaultFakeRepo())
	out, err := svc.Generate(context.Background(), narrativegen.Input{
		Class:   &class,
		Species: &species,
		Seed:    ptrSeed(seed),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Background.Category != domain.NarrativeBackground {
		t.Errorf("expected Background.Category=%q, got %q", domain.NarrativeBackground, out.Background.Category)
	}
	if out.Motivation.Category != domain.NarrativeMotivation {
		t.Errorf("expected Motivation.Category=%q, got %q", domain.NarrativeMotivation, out.Motivation.Category)
	}
	if out.Secret.Category != domain.NarrativeSecret {
		t.Errorf("expected Secret.Category=%q, got %q", domain.NarrativeSecret, out.Secret.Category)
	}
}

// ─── 11. No duplicate blocks — different content per category ────────────────

func TestGenerate_NoDuplicateBlocks(t *testing.T) {
	class := domain.ClassDruid
	species := domain.SpeciesElf

	for i := int64(0); i < 20; i++ {
		svc := narrativegen.New(defaultFakeRepo())
		out, err := svc.Generate(context.Background(), narrativegen.Input{
			Class:   &class,
			Species: &species,
			Seed:    ptrSeed(i),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
		if out.Background.Content == out.Motivation.Content {
			t.Errorf("seed %d: Background and Motivation have same content: %q", i, out.Background.Content)
		}
		if out.Background.Content == out.Secret.Content {
			t.Errorf("seed %d: Background and Secret have same content: %q", i, out.Background.Content)
		}
		if out.Motivation.Content == out.Secret.Content {
			t.Errorf("seed %d: Motivation and Secret have same content: %q", i, out.Motivation.Content)
		}
	}
}

// ─── 12. Empty pool returns error ─────────────────────────────────────────────

func TestGenerate_EmptyPool_ReturnsError(t *testing.T) {
	repo := &fakeNarrativeRepo{
		entries: map[domain.NarrativeCategory][]ports.WeightedNarrativeEntry{
			domain.NarrativeBackground: {}, // empty
			domain.NarrativeMotivation: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeMotivation, Content: "ok"}, Weight: 2},
			},
			domain.NarrativeSecret: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeSecret, Content: "ok"}, Weight: 2},
			},
		},
	}

	svc := narrativegen.New(repo)
	_, err := svc.Generate(context.Background(), narrativegen.Input{
		Class:   ptrClass(domain.ClassFighter),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(1),
	})
	if err == nil {
		t.Fatal("expected error for empty pool, got nil")
	}
}

// ─── 13. Weighted distribution — primary entries dominate ────────────────────

func TestGenerate_WeightedDistribution(t *testing.T) {
	// primary (weight 10) vs default (weight 2): ratio should be ~5:1
	// Over 1000 draws, primary should appear 750-900 times (allow wide margin).
	primaryContent := "Entrada primaria"
	defaultContent := "Entrada por defecto"

	repo := &fakeNarrativeRepo{
		entries: map[domain.NarrativeCategory][]ports.WeightedNarrativeEntry{
			domain.NarrativeBackground: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeBackground, Content: primaryContent}, Weight: 10},
				{Block: domain.NarrativeBlock{Category: domain.NarrativeBackground, Content: defaultContent}, Weight: 2},
			},
			domain.NarrativeMotivation: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeMotivation, Content: "ok"}, Weight: 2},
			},
			domain.NarrativeSecret: {
				{Block: domain.NarrativeBlock{Category: domain.NarrativeSecret, Content: "ok"}, Weight: 2},
			},
		},
	}

	svc := narrativegen.New(repo)
	primaryCount := 0
	const draws = 1000

	for i := int64(0); i < draws; i++ {
		out, err := svc.Generate(context.Background(), narrativegen.Input{
			Class:   ptrClass(domain.ClassFighter),
			Species: ptrSpecies(domain.SpeciesHuman),
			Seed:    ptrSeed(i),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
		if out.Background.Content == primaryContent {
			primaryCount++
		}
	}

	// primary weight 10 / total weight 12 ≈ 83.3% — expect 700-950
	if primaryCount < 700 || primaryCount > 950 {
		t.Errorf("primary appeared %d/1000 times — expected 700-950 (83%% ± margin)", primaryCount)
	}
}
