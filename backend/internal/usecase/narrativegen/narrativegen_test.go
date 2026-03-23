package narrativegen_test

import (
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/narrativegen"
)

// ptr helpers

func ptrClass(c domain.Class) *domain.Class     { return &c }
func ptrSpecies(s domain.Species) *domain.Species { return &s }
func ptrSeed(v int64) *int64                     { return &v }

// ─── 1. Reproducibility ──────────────────────────────────────────────────────

func TestGenerate_Reproducibility(t *testing.T) {
	seed := int64(42)
	class := domain.ClassFighter
	species := domain.SpeciesHuman

	in := narrativegen.Input{
		Class:   ptrClass(class),
		Species: ptrSpecies(species),
		Seed:    ptrSeed(seed),
	}

	out1, err1 := narrativegen.Generate(in)
	out2, err2 := narrativegen.Generate(in)

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
	seed := int64(99)
	out, err := narrativegen.Generate(narrativegen.Input{
		Class:   ptrClass(domain.ClassWizard),
		Species: ptrSpecies(domain.SpeciesElf),
		Seed:    ptrSeed(seed),
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

// ─── 3. Tag filtering — class-specific template must not appear for other class ─

// This test verifies that the tag filter works: we assert that if we look at
// the pool with a class that has specific templates, filtering works correctly.
// We test this indirectly: generate many times without seed for a Barbarian
// and verify we never get a Wizard-only template (we check the tags of the
// returned block are compatible with Barbarian).
func TestGenerate_TagFiltering_ClassSpecific(t *testing.T) {
	barbarian := domain.ClassBarbarian
	human := domain.SpeciesHuman

	for i := int64(0); i < 50; i++ {
		out, err := narrativegen.Generate(narrativegen.Input{
			Class:   &barbarian,
			Species: &human,
			Seed:    ptrSeed(i),
		})
		if err != nil {
			t.Fatalf("seed %d: %v", i, err)
		}
		// The returned block's tags must contain "any", "barbarian", or "human"
		if !blockIsCompatible(out.Background, barbarian, human) {
			t.Errorf("seed %d: Background block incompatible with barbarian/human: tags=%v content=%q",
				i, out.Background.Tags, out.Background.Content)
		}
		if !blockIsCompatible(out.Motivation, barbarian, human) {
			t.Errorf("seed %d: Motivation block incompatible with barbarian/human: tags=%v content=%q",
				i, out.Motivation.Tags, out.Motivation.Content)
		}
		if !blockIsCompatible(out.Secret, barbarian, human) {
			t.Errorf("seed %d: Secret block incompatible with barbarian/human: tags=%v content=%q",
				i, out.Secret.Tags, out.Secret.Content)
		}
	}
}

// blockIsCompatible mirrors the internal isCompatible logic.
func blockIsCompatible(b domain.NarrativeBlock, class domain.Class, species domain.Species) bool {
	for _, tag := range b.Tags {
		if tag == "any" || tag == string(class) || tag == string(species) {
			return true
		}
	}
	return false
}

// ─── 4. Universal tags ("any") compatible with every class ───────────────────

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
			out, err := narrativegen.Generate(narrativegen.Input{
				Class:   &cls,
				Species: &species,
				Seed:    ptrSeed(seed),
			})
			if err != nil {
				t.Fatalf("class %q: %v", cls, err)
			}
			if !blockIsCompatible(out.Background, cls, species) {
				t.Errorf("Background not compatible with class %q: tags=%v", cls, out.Background.Tags)
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
			out, err := narrativegen.Generate(narrativegen.Input{
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
			out, err := narrativegen.Generate(narrativegen.Input{
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
		out, err := narrativegen.Generate(narrativegen.Input{
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

	out, err := narrativegen.Generate(narrativegen.Input{
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

	out, err := narrativegen.Generate(narrativegen.Input{
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

	out, err := narrativegen.Generate(narrativegen.Input{
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
	// Run across many seeds and verify all three blocks have distinct content.
	class := domain.ClassDruid
	species := domain.SpeciesElf

	for i := int64(0); i < 20; i++ {
		out, err := narrativegen.Generate(narrativegen.Input{
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
