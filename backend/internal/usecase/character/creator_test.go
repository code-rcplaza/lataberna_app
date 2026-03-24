package character_test

import (
	"context"
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
	"forge-rpg/internal/usecase/character"
	"forge-rpg/internal/usecase/namegen"
	"forge-rpg/internal/usecase/narrativegen"
)

// ptr helpers — keep tests readable.
func ptrString(s string) *string                 { return &s }
func ptrInt64(n int64) *int64                    { return &n }
func ptrClass(c domain.Class) *domain.Class      { return &c }
func ptrSpecies(s domain.Species) *domain.Species { return &s }

// modifierFor computes the expected modifier from a final stat value.
func modifierFor(stat int) int {
	if stat >= 10 {
		return (stat - 10) / 2
	}
	return (stat - 10 - 1) / 2
}

// ---------------------------------------------------------------------------
// fakeNarrativeRepo — minimal stub for Creator tests
// ---------------------------------------------------------------------------

type fakeNarrativeRepo struct{}

func (f *fakeNarrativeRepo) FindByCategory(
	ctx      context.Context,
	category domain.NarrativeCategory,
	class    domain.Class,
	species  domain.Species,
) ([]ports.WeightedNarrativeEntry, error) {
	contents := map[domain.NarrativeCategory][]string{
		domain.NarrativeBackground: {
			"Fondo A", "Fondo B", "Fondo C", "Fondo D", "Fondo E",
		},
		domain.NarrativeMotivation: {
			"Motivación A", "Motivación B", "Motivación C", "Motivación D", "Motivación E",
		},
		domain.NarrativeSecret: {
			"Secreto A", "Secreto B", "Secreto C", "Secreto D", "Secreto E",
		},
	}
	pool := contents[category]
	out := make([]ports.WeightedNarrativeEntry, len(pool))
	for i, c := range pool {
		out[i] = ports.WeightedNarrativeEntry{
			Block:  domain.NarrativeBlock{Category: category, Content: c},
			Weight: 2,
		}
	}
	return out, nil
}

func (f *fakeNarrativeRepo) Count(ctx context.Context) (int, error) { return 15, nil }

// ---------------------------------------------------------------------------
// fakeNameRepo — minimal stub for Creator tests
// ---------------------------------------------------------------------------

type fakeNameRepo struct{}

func (f *fakeNameRepo) FindBySpeciesGender(ctx context.Context, speciesKey, gender string) ([]string, error) {
	return f.FindByType(ctx, speciesKey, gender, "first_name")
}

func (f *fakeNameRepo) FindByType(ctx context.Context, speciesKey, gender, nameType string) ([]string, error) {
	return []string{
		"Aldric", "Brennan", "Cael", "Dorian", "Edric",
		"Faolan", "Gareth", "Hadwin", "Isidor", "Jareth",
		"Kiran", "Leoric", "Maddox", "Nolan", "Orwin",
		"Phelan", "Quinn", "Roderick", "Soren", "Theron",
		"Ulric", "Vance", "Wulfric", "Xander", "Yorick",
	}, nil
}

func (f *fakeNameRepo) Count(ctx context.Context) (int, error) { return 25, nil }

// newTestCreator returns a Creator with fake repos for unit testing.
func newTestCreator() *character.Creator {
	narrativeSvc := narrativegen.New(&fakeNarrativeRepo{})
	nameSvc := namegen.New(&fakeNameRepo{})
	return character.NewCreator(narrativeSvc, nameSvc)
}

// --- Test: zero inputs produce a valid character ---

func TestCreate_ZeroInputs_ReturnsValidCharacter(t *testing.T) {
	creator := newTestCreator()
	out, err := creator.Create(context.Background(), character.CreateInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil character")
	}
	if out.Name == "" {
		t.Error("Name should not be empty")
	}
	if out.Class == "" {
		t.Error("Class should not be empty")
	}
	if out.Species == "" {
		t.Error("Species should not be empty")
	}
	if out.Level != 1 {
		t.Errorf("Level should be 1, got %d", out.Level)
	}
	if out.Ruleset != domain.Ruleset5e {
		t.Errorf("Ruleset should be %q, got %q", domain.Ruleset5e, out.Ruleset)
	}
	if out.AbilityBonusSource != domain.AbilityBonusFromSpecies {
		t.Errorf("AbilityBonusSource should be %q, got %q", domain.AbilityBonusFromSpecies, out.AbilityBonusSource)
	}
	if out.Background.Content == "" {
		t.Error("Background.Content should not be empty")
	}
	if out.Motivation.Content == "" {
		t.Error("Motivation.Content should not be empty")
	}
	if out.Secret.Content == "" {
		t.Error("Secret.Content should not be empty")
	}
}

// --- Test: same seed → identical character ---

func TestCreate_SameSeed_ReproducibleResult(t *testing.T) {
	creator := newTestCreator()
	seed := int64(42)
	a, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(seed)})
	if err != nil {
		t.Fatalf("first Create: %v", err)
	}
	b, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(seed)})
	if err != nil {
		t.Fatalf("second Create: %v", err)
	}

	if a.Name != b.Name {
		t.Errorf("Name mismatch: %q vs %q", a.Name, b.Name)
	}
	if a.Class != b.Class {
		t.Errorf("Class mismatch: %q vs %q", a.Class, b.Class)
	}
	if a.Species != b.Species {
		t.Errorf("Species mismatch: %q vs %q", a.Species, b.Species)
	}
	if a.BaseStats != b.BaseStats {
		t.Errorf("BaseStats mismatch: %+v vs %+v", a.BaseStats, b.BaseStats)
	}
	if a.FinalStats != b.FinalStats {
		t.Errorf("FinalStats mismatch: %+v vs %+v", a.FinalStats, b.FinalStats)
	}
	if a.Modifiers != b.Modifiers {
		t.Errorf("Modifiers mismatch: %+v vs %+v", a.Modifiers, b.Modifiers)
	}
	if a.Background.Content != b.Background.Content {
		t.Errorf("Background mismatch: %q vs %q", a.Background.Content, b.Background.Content)
	}
	if a.Motivation.Content != b.Motivation.Content {
		t.Errorf("Motivation mismatch: %q vs %q", a.Motivation.Content, b.Motivation.Content)
	}
	if a.Secret.Content != b.Secret.Content {
		t.Errorf("Secret mismatch: %q vs %q", a.Secret.Content, b.Secret.Content)
	}
}

// --- Test: provided class is preserved ---

func TestCreate_ProvidedClass_IsPreserved(t *testing.T) {
	creator := newTestCreator()
	class := domain.ClassFighter
	out, err := creator.Create(context.Background(), character.CreateInput{Class: ptrClass(class)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Class != class {
		t.Errorf("expected class %q, got %q", class, out.Class)
	}
}

// --- Test: provided species is preserved ---

func TestCreate_ProvidedSpecies_IsPreserved(t *testing.T) {
	creator := newTestCreator()
	species := domain.SpeciesElf
	out, err := creator.Create(context.Background(), character.CreateInput{Species: ptrSpecies(species)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Species != species {
		t.Errorf("expected species %q, got %q", species, out.Species)
	}
}

// --- Test: provided name is preserved ---

func TestCreate_ProvidedName_IsPreserved(t *testing.T) {
	creator := newTestCreator()
	name := "Aldric"
	out, err := creator.Create(context.Background(), character.CreateInput{Name: ptrString(name)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name != name {
		t.Errorf("expected name %q, got %q", name, out.Name)
	}
}

// --- Test: all 13 classes generate without error ---

func TestCreate_AllClasses_GenerateWithoutError(t *testing.T) {
	classes := []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric, domain.ClassDruid,
		domain.ClassFighter, domain.ClassMonk, domain.ClassPaladin, domain.ClassRanger,
		domain.ClassRogue, domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}
	for _, class := range classes {
		t.Run(string(class), func(t *testing.T) {
			creator := newTestCreator()
			out, err := creator.Create(context.Background(), character.CreateInput{Class: ptrClass(class)})
			if err != nil {
				t.Fatalf("class %q: expected no error, got %v", class, err)
			}
			if out == nil {
				t.Fatalf("class %q: expected non-nil character", class)
			}
		})
	}
}

// --- Test: modifiers match finalStats ---

func TestCreate_Modifiers_MatchFinalStats(t *testing.T) {
	creator := newTestCreator()
	seed := int64(99)
	out, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(seed)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	checks := []struct {
		name string
		got  int
		stat int
	}{
		{"STR", out.Modifiers.STR, out.FinalStats.STR},
		{"DEX", out.Modifiers.DEX, out.FinalStats.DEX},
		{"CON", out.Modifiers.CON, out.FinalStats.CON},
		{"INT", out.Modifiers.INT, out.FinalStats.INT},
		{"WIS", out.Modifiers.WIS, out.FinalStats.WIS},
		{"CHA", out.Modifiers.CHA, out.FinalStats.CHA},
	}
	for _, c := range checks {
		want := modifierFor(c.stat)
		if c.got != want {
			t.Errorf("%s modifier: got %d, want %d (finalStat=%d)", c.name, c.got, want, c.stat)
		}
	}
}

// --- Test: Regenerate with nil character returns error ---

func TestRegenerate_NilCharacter_ReturnsError(t *testing.T) {
	creator := newTestCreator()
	_, err := creator.Regenerate(context.Background(), character.RegenerateInput{
		Character: nil,
		Locks:     domain.CharacterLocks{},
	})
	if err == nil {
		t.Fatal("expected error for nil character, got nil")
	}
}

// --- Test: locked name is preserved on Regenerate ---

func TestRegenerate_NameLocked_NameUnchanged(t *testing.T) {
	creator := newTestCreator()
	original, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(1)})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	regen, err := creator.Regenerate(context.Background(), character.RegenerateInput{
		Character: original,
		Locks:     domain.CharacterLocks{Name: true},
		Seed:      ptrInt64(2),
	})
	if err != nil {
		t.Fatalf("Regenerate: %v", err)
	}
	if regen.Name != original.Name {
		t.Errorf("expected locked name %q, got %q", original.Name, regen.Name)
	}
}

// --- Test: unlocked name changes on Regenerate with different seed ---

func TestRegenerate_NameUnlocked_NameChanges(t *testing.T) {
	creator := newTestCreator()
	original, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(111)})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	changed := false
	for seed := int64(9999); seed < 10100; seed++ {
		regen, err := creator.Regenerate(context.Background(), character.RegenerateInput{
			Character: original,
			Locks:     domain.CharacterLocks{Name: false},
			Seed:      ptrInt64(seed),
		})
		if err != nil {
			t.Fatalf("Regenerate seed %d: %v", seed, err)
		}
		if regen.Name != original.Name {
			changed = true
			break
		}
	}
	if !changed {
		t.Error("expected name to change after unlocked regeneration with different seed, but it never did")
	}
}

// --- Test: stats locked → stats unchanged on Regenerate ---

func TestRegenerate_StatsLocked_StatsUnchanged(t *testing.T) {
	creator := newTestCreator()
	original, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(200)})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	regen, err := creator.Regenerate(context.Background(), character.RegenerateInput{
		Character: original,
		Locks:     domain.CharacterLocks{Stats: true},
		Seed:      ptrInt64(300),
	})
	if err != nil {
		t.Fatalf("Regenerate: %v", err)
	}
	if regen.BaseStats != original.BaseStats {
		t.Errorf("BaseStats changed: got %+v, want %+v", regen.BaseStats, original.BaseStats)
	}
	if regen.FinalStats != original.FinalStats {
		t.Errorf("FinalStats changed: got %+v, want %+v", regen.FinalStats, original.FinalStats)
	}
	if regen.Modifiers != original.Modifiers {
		t.Errorf("Modifiers changed: got %+v, want %+v", regen.Modifiers, original.Modifiers)
	}
	if regen.Derived != original.Derived {
		t.Errorf("Derived changed: got %+v, want %+v", regen.Derived, original.Derived)
	}
}

// --- Test: narrative blocks locked individually ---

func TestRegenerate_BackgroundLocked_OnlyBackgroundUnchanged(t *testing.T) {
	creator := newTestCreator()
	original, err := creator.Create(context.Background(), character.CreateInput{Seed: ptrInt64(500)})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	regenerated := false
	for seed := int64(8000); seed < 8100; seed++ {
		regen, err := creator.Regenerate(context.Background(), character.RegenerateInput{
			Character: original,
			Locks: domain.CharacterLocks{
				Background: true,
				Motivation: false,
				Secret:     false,
			},
			Seed: ptrInt64(seed),
		})
		if err != nil {
			t.Fatalf("Regenerate seed %d: %v", seed, err)
		}
		if regen.Background.Content != original.Background.Content {
			t.Errorf("Background changed but was locked: got %q, want %q",
				regen.Background.Content, original.Background.Content)
		}
		if regen.Motivation.Content != original.Motivation.Content ||
			regen.Secret.Content != original.Secret.Content {
			regenerated = true
			break
		}
	}
	if !regenerated {
		t.Error("expected Motivation or Secret to change after unlocked regeneration, but neither did")
	}
}

// --- Test: character has non-empty ID ---

func TestCreate_CharacterHasNonEmptyID(t *testing.T) {
	creator := newTestCreator()
	out, err := creator.Create(context.Background(), character.CreateInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.ID == "" {
		t.Error("expected non-empty ID")
	}
}

// --- Test: character has non-zero timestamps ---

func TestCreate_CharacterHasTimestamps(t *testing.T) {
	creator := newTestCreator()
	out, err := creator.Create(context.Background(), character.CreateInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if out.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

// --- Test: Locks field is all-false on Create ---

func TestCreate_LocksAllFalseOnCreation(t *testing.T) {
	creator := newTestCreator()
	out, err := creator.Create(context.Background(), character.CreateInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Locks != (domain.CharacterLocks{}) {
		t.Errorf("expected all-false locks, got %+v", out.Locks)
	}
}

// --- Test: provided gender produces coherent name ---

func TestCreate_ProvidedGender_CoherentName(t *testing.T) {
	creator := newTestCreator()
	gender := namegen.GenderFemale
	out, err := creator.Create(context.Background(), character.CreateInput{
		Gender: &gender,
		Seed:   ptrInt64(77),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Name == "" {
		t.Error("Name should not be empty when gender is provided")
	}
}
