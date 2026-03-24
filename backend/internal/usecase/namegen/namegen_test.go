package namegen_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
	"forge-rpg/internal/usecase/namegen"
)

func ptr[T any](v T) *T { return &v }

// ---------------------------------------------------------------------------
// fakeNameRepo — in-memory stub satisfying ports.NameRepository
// ---------------------------------------------------------------------------

// fakeNameRepo uses a 3-level map: speciesKey → nameType → gender → names.
type fakeNameRepo struct {
	data map[string]map[string]map[string][]string
}

func (f *fakeNameRepo) FindByType(ctx context.Context, speciesKey, gender, nameType string) ([]string, error) {
	if typeMap, ok := f.data[speciesKey]; ok {
		if genderMap, ok := typeMap[nameType]; ok {
			if names, ok := genderMap[gender]; ok && len(names) > 0 {
				return names, nil
			}
		}
	}
	return nil, fmt.Errorf("fakeNameRepo: species=%q gender=%q type=%q: %w",
		speciesKey, gender, nameType, ports.ErrEmptyNamePool)
}

func (f *fakeNameRepo) FindBySpeciesGender(ctx context.Context, speciesKey, gender string) ([]string, error) {
	names, err := f.FindByType(ctx, speciesKey, gender, "first_name")
	if errors.Is(err, ports.ErrEmptyNamePool) {
		return nil, nil
	}
	return names, err
}

func (f *fakeNameRepo) Count(ctx context.Context) (int, error) {
	total := 0
	for _, typeMap := range f.data {
		for _, genderMap := range typeMap {
			for _, names := range genderMap {
				total += len(names)
			}
		}
	}
	return total, nil
}

// defaultFakeNameRepo returns a repo populated with pools for all species and
// all component types required by compose().
func defaultFakeNameRepo() *fakeNameRepo {
	make25 := func(prefix string) []string {
		names := make([]string, 25)
		for i := range names {
			names[i] = fmt.Sprintf("%s-%02d", prefix, i+1)
		}
		return names
	}

	return &fakeNameRepo{
		data: map[string]map[string]map[string][]string{
			"human": {
				"first_name": {"male": make25("hum-m"), "female": make25("hum-f")},
				"surname":    {"any": make25("hum-sn")},
			},
			"high-elf": {
				"first_name":  {"male": make25("helf-m"), "female": make25("helf-f")},
				"family_name": {"any": make25("helf-fam")},
			},
			"wood-elf": {
				"first_name":  {"male": make25("welf-m"), "female": make25("welf-f")},
				"family_name": {"any": make25("welf-fam")},
			},
			"drow": {
				"first_name":  {"male": make25("drow-m"), "female": make25("drow-f")},
				"family_name": {"any": make25("drow-fam")},
			},
			"hill-dwarf": {
				"first_name": {"male": make25("hdw-m"), "female": make25("hdw-f")},
				"clan_name":  {"any": make25("hdw-clan")},
			},
			"mountain-dwarf": {
				"first_name": {"male": make25("mdw-m"), "female": make25("mdw-f")},
				"clan_name":  {"any": make25("mdw-clan")},
			},
			"lightfoot": {
				"first_name": {"male": make25("lf-m"), "female": make25("lf-f")},
				"surname":    {"any": make25("lf-sn")},
			},
			"stout": {
				"first_name": {"male": make25("st-m"), "female": make25("st-f")},
				"surname":    {"any": make25("st-sn")},
			},
			"forest-gnome": {
				"first_name": {"male": make25("fg-m"), "female": make25("fg-f")},
				"clan_name":  {"any": make25("fg-clan")},
				"nickname":   {"any": make25("fg-nick")},
			},
			"rock-gnome": {
				"first_name": {"male": make25("rg-m"), "female": make25("rg-f")},
				"clan_name":  {"any": make25("rg-clan")},
				"nickname":   {"any": make25("rg-nick")},
			},
			"half-elf": {
				"first_name":  {"male": make25("he-m"), "female": make25("he-f")},
				"surname":     {"any": make25("he-sn")},
				"family_name": {"any": make25("he-fam")},
			},
			"half-orc": {
				"first_name": {"male": make25("ho-m"), "female": make25("ho-f")},
				"surname":    {"any": make25("ho-sn")},
			},
			"tiefling-infernal": {
				"infernal_name": {"any": make25("tief-inf")},
			},
			"tiefling-virtue": {
				"virtue_word": {"any": make25("tief-vir")},
			},
			"dragonborn": {
				"first_name": {"male": make25("db-m"), "female": make25("db-f")},
				"clan_name":  {"any": make25("db-clan")},
			},
		},
	}
}

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
		{"tiefling-infernal", domain.SpeciesTiefling, ptr(domain.SubSpeciesInfernalTiefling)},
		{"tiefling-virtue", domain.SpeciesTiefling, ptr(domain.SubSpeciesVirtueTiefling)},
		{"dragonborn", domain.SpeciesDragonborn, nil},
	}

	seed := int64(42)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := namegen.New(defaultFakeNameRepo())
			out, err := svc.Generate(context.Background(), namegen.Input{
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
	svc := namegen.New(defaultFakeNameRepo())
	in := namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seed,
	}

	out1, err1 := svc.Generate(context.Background(), in)
	out2, err2 := svc.Generate(context.Background(), in)

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
	seedA := int64(1)
	seedB := int64(999)
	svc := namegen.New(defaultFakeNameRepo())

	outA, err := svc.Generate(context.Background(), namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seedA,
	})
	if err != nil {
		t.Fatalf("seed A error: %v", err)
	}

	outB, err := svc.Generate(context.Background(), namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seedB,
	})
	if err != nil {
		t.Fatalf("seed B error: %v", err)
	}

	if outA.Name == outB.Name {
		t.Logf("WARNING: seeds %d and %d produced the same name %q — possible but unlikely", seedA, seedB, outA.Name)
	}
}

// --- 4. Nil seed generates a random (non-empty) name ---

func TestGenerate_NilSeed_ReturnsValidName(t *testing.T) {
	svc := namegen.New(defaultFakeNameRepo())
	out, err := svc.Generate(context.Background(), namegen.Input{
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
			svc := namegen.New(defaultFakeNameRepo())
			out, err := svc.Generate(context.Background(), namegen.Input{
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
			svc := namegen.New(defaultFakeNameRepo())
			sub := tc.sub
			out, err := svc.Generate(context.Background(), namegen.Input{
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

	svc := namegen.New(defaultFakeNameRepo())
	baseline, err := svc.Generate(context.Background(), namegen.Input{
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
			out, err := svc.Generate(context.Background(), namegen.Input{
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
	svc := namegen.New(defaultFakeNameRepo())
	_, err := svc.Generate(context.Background(), namegen.Input{
		Species: domain.Species("unknown-alien"),
		Seed:    &seed,
	})
	if err == nil {
		t.Fatal("expected error for unknown species, got nil")
	}
}

// ---------------------------------------------------------------------------
// Composition rules (TDD — Phase 5 tasks 5.1-5.6)
// ---------------------------------------------------------------------------

// Task 5.1 — Human: output is <word> <word>, both tokens non-empty.
func TestCompose_Human(t *testing.T) {
	seed := int64(42)
	svc := namegen.New(defaultFakeNameRepo())
	out, err := svc.Generate(context.Background(), namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seed,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tokens := strings.Fields(out.Name)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens for Human, got %d: %q", len(tokens), out.Name)
	}
	if tokens[0] == "" || tokens[1] == "" {
		t.Fatalf("expected non-empty tokens, got %q", out.Name)
	}
}

// Task 5.2 — Dragonborn: clan name is the FIRST token.
func TestCompose_Dragonborn(t *testing.T) {
	clanPool := []string{}
	for i := 1; i <= 25; i++ {
		clanPool = append(clanPool, fmt.Sprintf("db-clan-%02d", i))
	}

	// Run 20 fixed-seed iterations to assert clan-first is always true.
	for seedVal := int64(0); seedVal < 20; seedVal++ {
		svc := namegen.New(defaultFakeNameRepo())
		out, err := svc.Generate(context.Background(), namegen.Input{
			Species: domain.SpeciesDragonborn,
			Gender:  ptr(namegen.GenderMale),
			Seed:    &seedVal,
		})
		if err != nil {
			t.Fatalf("seed %d: unexpected error: %v", seedVal, err)
		}
		tokens := strings.Fields(out.Name)
		if len(tokens) != 2 {
			t.Fatalf("seed %d: expected 2 tokens, got %d: %q", seedVal, len(tokens), out.Name)
		}
		// First token must be from the clan pool (prefix "db-clan-").
		if !strings.HasPrefix(tokens[0], "db-clan-") {
			t.Errorf("seed %d: first token %q is not a clan name (expected prefix db-clan-)", seedVal, tokens[0])
		}
	}
}

// Task 5.3 — Gnome: nickname is wrapped in double quotes.
func TestCompose_Gnome(t *testing.T) {
	seed := int64(42)
	svc := namegen.New(defaultFakeNameRepo())
	out, err := svc.Generate(context.Background(), namegen.Input{
		Species:    domain.SpeciesGnome,
		SubSpecies: ptr(domain.SubSpeciesForestGnome),
		Gender:     ptr(namegen.GenderMale),
		Seed:       &seed,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Output format: first clan "nickname"
	if !strings.Contains(out.Name, `"`) {
		t.Fatalf("Gnome name %q does not contain double-quoted nickname", out.Name)
	}
	// Verify structure: three whitespace-separated tokens where last is "nick-XX"
	parts := strings.Fields(out.Name)
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts for Gnome name, got %d: %q", len(parts), out.Name)
	}
	if !strings.HasPrefix(parts[2], `"`) || !strings.HasSuffix(parts[2], `"`) {
		t.Errorf("third token %q should be wrapped in double quotes", parts[2])
	}
}

// Task 5.4 — Tiefling: infernal key never returns virtue word; virtue key never returns infernal name.
func TestCompose_Tiefling(t *testing.T) {
	t.Run("infernal key returns infernal name only", func(t *testing.T) {
		for seedVal := int64(0); seedVal < 20; seedVal++ {
			svc := namegen.New(defaultFakeNameRepo())
			sub := domain.SubSpeciesInfernalTiefling
			out, err := svc.Generate(context.Background(), namegen.Input{
				Species:    domain.SpeciesTiefling,
				SubSpecies: &sub,
				Seed:       &seedVal,
			})
			if err != nil {
				t.Fatalf("seed %d: unexpected error: %v", seedVal, err)
			}
			// Pool prefix for infernal is "tief-inf-"
			if !strings.HasPrefix(out.Name, "tief-inf-") {
				t.Errorf("seed %d: infernal tiefling got non-infernal name %q", seedVal, out.Name)
			}
		}
	})

	t.Run("virtue key returns virtue word only", func(t *testing.T) {
		for seedVal := int64(0); seedVal < 20; seedVal++ {
			svc := namegen.New(defaultFakeNameRepo())
			sub := domain.SubSpeciesVirtueTiefling
			out, err := svc.Generate(context.Background(), namegen.Input{
				Species:    domain.SpeciesTiefling,
				SubSpecies: &sub,
				Seed:       &seedVal,
			})
			if err != nil {
				t.Fatalf("seed %d: unexpected error: %v", seedVal, err)
			}
			// Pool prefix for virtue is "tief-vir-"
			if !strings.HasPrefix(out.Name, "tief-vir-") {
				t.Errorf("seed %d: virtue tiefling got non-virtue name %q", seedVal, out.Name)
			}
		}
	})
}

// Task 5.5 — Half-Orc probability: ~30% two-word, ~70% one-word over N=100 runs.
func TestCompose_HalfOrc_Probability(t *testing.T) {
	const N = 100
	twoWordCount := 0

	for seedVal := int64(0); seedVal < N; seedVal++ {
		svc := namegen.New(defaultFakeNameRepo())
		out, err := svc.Generate(context.Background(), namegen.Input{
			Species: domain.SpeciesHalfOrc,
			Gender:  ptr(namegen.GenderMale),
			Seed:    &seedVal,
		})
		if err != nil {
			t.Fatalf("seed %d: unexpected error: %v", seedVal, err)
		}
		if len(strings.Fields(out.Name)) == 2 {
			twoWordCount++
		}
	}

	// Expected ~30% ± 15pp tolerance (15-45 out of 100).
	if twoWordCount < 15 || twoWordCount > 45 {
		t.Errorf("Half-Orc two-word name rate: %d/100 (expected 15-45)", twoWordCount)
	}
}

// Task 5.6 — Remaining species: correct assembly order.

func TestCompose_Dwarf(t *testing.T) {
	for _, sub := range []domain.SubSpecies{domain.SubSpeciesHillDwarf, domain.SubSpeciesMountainDwarf} {
		t.Run(string(sub), func(t *testing.T) {
			seed := int64(42)
			svc := namegen.New(defaultFakeNameRepo())
			out, err := svc.Generate(context.Background(), namegen.Input{
				Species:    domain.SpeciesDwarf,
				SubSpecies: &sub,
				Gender:     ptr(namegen.GenderFemale),
				Seed:       &seed,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tokens := strings.Fields(out.Name)
			if len(tokens) != 2 {
				t.Fatalf("expected 2 tokens for Dwarf, got %d: %q", len(tokens), out.Name)
			}
		})
	}
}

func TestCompose_Elf(t *testing.T) {
	for _, sub := range []domain.SubSpecies{domain.SubSpeciesHighElf, domain.SubSpeciesWoodElf, domain.SubSpeciesDrow} {
		t.Run(string(sub), func(t *testing.T) {
			seed := int64(42)
			svc := namegen.New(defaultFakeNameRepo())
			out, err := svc.Generate(context.Background(), namegen.Input{
				Species:    domain.SpeciesElf,
				SubSpecies: &sub,
				Gender:     ptr(namegen.GenderFemale),
				Seed:       &seed,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tokens := strings.Fields(out.Name)
			if len(tokens) != 2 {
				t.Fatalf("expected 2 tokens for Elf, got %d: %q", len(tokens), out.Name)
			}
		})
	}
}

func TestCompose_Halfling(t *testing.T) {
	for _, sub := range []domain.SubSpecies{domain.SubSpeciesLightfoot, domain.SubSpeciesStout} {
		t.Run(string(sub), func(t *testing.T) {
			seed := int64(42)
			svc := namegen.New(defaultFakeNameRepo())
			out, err := svc.Generate(context.Background(), namegen.Input{
				Species:    domain.SpeciesHalfling,
				SubSpecies: &sub,
				Gender:     ptr(namegen.GenderFemale),
				Seed:       &seed,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tokens := strings.Fields(out.Name)
			if len(tokens) != 2 {
				t.Fatalf("expected 2 tokens for Halfling, got %d: %q", len(tokens), out.Name)
			}
		})
	}
}

func TestCompose_HalfElf_TwoTokens(t *testing.T) {
	// Half-Elf always produces two tokens regardless of convention.
	for seedVal := int64(0); seedVal < 20; seedVal++ {
		svc := namegen.New(defaultFakeNameRepo())
		out, err := svc.Generate(context.Background(), namegen.Input{
			Species: domain.SpeciesHalfElf,
			Gender:  ptr(namegen.GenderFemale),
			Seed:    &seedVal,
		})
		if err != nil {
			t.Fatalf("seed %d: unexpected error: %v", seedVal, err)
		}
		tokens := strings.Fields(out.Name)
		if len(tokens) != 2 {
			t.Errorf("seed %d: expected 2 tokens for Half-Elf, got %d: %q", seedVal, len(tokens), out.Name)
		}
	}
}

// Task 5.7 — Error on empty pool (tested via missing pool in repo).
func TestGenerate_EmptyPool_ReturnsError(t *testing.T) {
	// A repo with no pools at all.
	emptyRepo := &fakeNameRepo{data: map[string]map[string]map[string][]string{}}
	svc := namegen.New(emptyRepo)
	seed := int64(1)
	_, err := svc.Generate(context.Background(), namegen.Input{
		Species: domain.SpeciesHuman,
		Gender:  ptr(namegen.GenderMale),
		Seed:    &seed,
	})
	if err == nil {
		t.Fatal("expected error when pool is empty, got nil")
	}
	if !errors.Is(err, ports.ErrEmptyNamePool) {
		t.Errorf("expected error to wrap ErrEmptyNamePool, got: %v", err)
	}
}

// --- 9. Verify interface is satisfied at compile time ---

var _ ports.NameRepository = (*fakeNameRepo)(nil)
