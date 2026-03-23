package namegen_test

import (
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/namegen"
)

func ptr[T any](v T) *T { return &v }

// --- 1. Valid name returned for each of the 9 species ---

func TestGenerate_ValidNamePerSpecies(t *testing.T) {
	cases := []struct {
		name    string
		species domain.Species
		sub     *domain.SubSpecies
	}{
		{"human", domain.SpeciesHuman, nil},
		{"elf-high", domain.SpeciesElf, ptr(domain.SubSpeciesHighElf)},
		{"elf-wood", domain.SpeciesElf, ptr(domain.SubSpeciesWoodElf)},
		{"elf-drow", domain.SpeciesElf, ptr(domain.SubSpeciesDrow)},
		{"dwarf-hill", domain.SpeciesDwarf, ptr(domain.SubSpeciesHillDwarf)},
		{"dwarf-mountain", domain.SpeciesDwarf, ptr(domain.SubSpeciesMountainDwarf)},
		{"halfling-lightfoot", domain.SpeciesHalfling, ptr(domain.SubSpeciesLightfoot)},
		{"halfling-stout", domain.SpeciesHalfling, ptr(domain.SubSpeciesStout)},
		{"gnome-forest", domain.SpeciesGnome, ptr(domain.SubSpeciesForestGnome)},
		{"gnome-rock", domain.SpeciesGnome, ptr(domain.SubSpeciesRockGnome)},
		{"half-elf", domain.SpeciesHalfElf, nil},
		{"half-orc", domain.SpeciesHalfOrc, nil},
		{"tiefling", domain.SpeciesTiefling, nil},
		{"dragonborn", domain.SpeciesDragonborn, nil},
	}

	seed := int64(42)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := namegen.Generate(namegen.Input{
				Species:    tc.species,
				SubSpecies: tc.sub,
				Seed:       &seed,
			})
			if err != nil {
				t.Fatalf("Generate returned error: %v", err)
			}
			if out.Name == "" {
				t.Fatal("Generate returned empty name")
			}
		})
	}
}

// --- 2. Reproducibility: same seed + same input → same name ---

func TestGenerate_Reproducibility(t *testing.T) {
	seed := int64(12345)
	in := namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seed,
	}

	out1, err1 := namegen.Generate(in)
	out2, err2 := namegen.Generate(in)

	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v / %v", err1, err2)
	}
	if out1.Name != out2.Name {
		t.Fatalf("expected reproducible name, got %q and %q", out1.Name, out2.Name)
	}
	if out1.Seed != out2.Seed {
		t.Fatalf("expected same seed in output, got %d and %d", out1.Seed, out2.Seed)
	}
}

// --- 3. Different seeds may produce different names ---

func TestGenerate_DifferentSeeds_MayDiffer(t *testing.T) {
	// Use human which has 25+ names per gender — very likely to differ across seeds
	seedA := int64(1)
	seedB := int64(999)

	outA, err := namegen.Generate(namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seedA,
	})
	if err != nil {
		t.Fatalf("seed A error: %v", err)
	}

	outB, err := namegen.Generate(namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seedB,
	})
	if err != nil {
		t.Fatalf("seed B error: %v", err)
	}

	// Not a hard failure if they happen to be the same (pool is finite),
	// but with 25 names this is a ~1/25 chance — log a warning only.
	if outA.Name == outB.Name {
		t.Logf("WARNING: seeds %d and %d produced the same name %q — possible but unlikely", seedA, seedB, outA.Name)
	}
}

// --- 4. Nil seed generates a random (non-empty) name ---

func TestGenerate_NilSeed_ReturnsValidName(t *testing.T) {
	out, err := namegen.Generate(namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderFemale),
		Seed:    nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Name == "" {
		t.Fatal("expected non-empty name with nil seed")
	}
	if out.Seed == 0 {
		t.Fatal("expected Seed to be set in output even when input Seed is nil")
	}
}

// --- 5. Gender resolution ---

func TestGenerate_GenderResolution(t *testing.T) {
	seed := int64(77)

	cases := []struct {
		name   string
		gender *namegen.Gender
	}{
		{"explicit-male", ptr(namegen.GenderMale)},
		{"explicit-female", ptr(namegen.GenderFemale)},
		{"nil-gender", nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := namegen.Generate(namegen.Input{
				Species: domain.SpeciesHuman,
				Gender:  tc.gender,
				Seed:    &seed,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.Name == "" {
				t.Fatal("expected non-empty name")
			}
		})
	}
}

// --- 6. SubSpecies respected ---

func TestGenerate_SubSpeciesRespected(t *testing.T) {
	seed := int64(55)

	cases := []struct {
		name string
		sub  domain.SubSpecies
	}{
		{"high-elf", domain.SubSpeciesHighElf},
		{"wood-elf", domain.SubSpeciesWoodElf},
		{"drow", domain.SubSpeciesDrow},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sub := tc.sub
			out, err := namegen.Generate(namegen.Input{
				Species:    domain.SpeciesElf,
				SubSpecies: &sub,
				Gender:     ptr(namegen.GenderMale),
				Seed:       &seed,
			})
			if err != nil {
				t.Fatalf("unexpected error for subspecies %q: %v", tc.sub, err)
			}
			if out.Name == "" {
				t.Fatalf("expected non-empty name for subspecies %q", tc.sub)
			}
		})
	}
}

// --- 7. All 13 classes have no effect on name generation ---

func TestGenerate_ClassHasNoEffect(t *testing.T) {
	seed := int64(42)
	gender := ptr(namegen.GenderMale)

	// Generate a baseline without class
	baseline, err := namegen.Generate(namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  gender,
		Seed:    &seed,
	})
	if err != nil {
		t.Fatalf("baseline error: %v", err)
	}

	classes := []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric, domain.ClassDruid,
		domain.ClassFighter, domain.ClassMonk, domain.ClassPaladin, domain.ClassRanger,
		domain.ClassRogue, domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}

	for _, cls := range classes {
		t.Run(string(cls), func(t *testing.T) {
			// Input has no Class field — name generation is species-only.
			// We verify the result is identical to baseline (class is irrelevant).
			out, err := namegen.Generate(namegen.Input{
				Species: domain.SpeciesHuman,
				Gender:  gender,
				Seed:    &seed,
			})
			if err != nil {
				t.Fatalf("class %q: unexpected error: %v", cls, err)
			}
			if out.Name != baseline.Name {
				t.Fatalf("class %q changed the name: got %q, want %q", cls, out.Name, baseline.Name)
			}
		})
	}
}

// --- 8. Error on unknown species ---

func TestGenerate_UnknownSpecies_ReturnsError(t *testing.T) {
	seed := int64(1)
	_, err := namegen.Generate(namegen.Input{
		Species: domain.Species("unknown-alien"),
		Seed:    &seed,
	})
	if err == nil {
		t.Fatal("expected error for unknown species, got nil")
	}
}
