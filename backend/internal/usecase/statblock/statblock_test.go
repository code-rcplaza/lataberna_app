package statblock_test

import (
	"testing"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/usecase/statblock"
)

// ptr helpers
func ptrClass(c domain.Class) *domain.Class       { return &c }
func ptrSpecies(s domain.Species) *domain.Species { return &s }
func ptrSub(s domain.SubSpecies) *domain.SubSpecies { return &s }
func ptrSeed(v int64) *int64                       { return &v }

// ─── 1. Reproducibility ──────────────────────────────────────────────────────

func TestGenerate_Reproducibility(t *testing.T) {
	seed := int64(42)
	class := domain.ClassFighter
	species := domain.SpeciesHuman

	in := statblock.Input{
		Class:   ptrClass(class),
		Species: ptrSpecies(species),
		Seed:    ptrSeed(seed),
	}

	out1, err1 := statblock.Generate(in)
	out2, err2 := statblock.Generate(in)

	if err1 != nil {
		t.Fatalf("first Generate: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("second Generate: %v", err2)
	}

	if out1.BaseStats != out2.BaseStats {
		t.Errorf("BaseStats differ: %+v vs %+v", out1.BaseStats, out2.BaseStats)
	}
	if out1.FinalStats != out2.FinalStats {
		t.Errorf("FinalStats differ: %+v vs %+v", out1.FinalStats, out2.FinalStats)
	}
	if out1.Modifiers != out2.Modifiers {
		t.Errorf("Modifiers differ: %+v vs %+v", out1.Modifiers, out2.Modifiers)
	}
	if out1.Derived != out2.Derived {
		t.Errorf("Derived differ: %+v vs %+v", out1.Derived, out2.Derived)
	}
	if out1.Seed != seed {
		t.Errorf("Seed: got %d, want %d", out1.Seed, seed)
	}
}

// ─── 2. Modifiers always from FinalStats ─────────────────────────────────────

func TestGenerate_ModifiersFromFinalStats(t *testing.T) {
	seed := int64(7)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassWizard),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}

	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Use CalculateModifier as the reference — it is independently tested in TestCalculateModifier.
	// This test verifies that out.Modifiers are derived from FinalStats, not BaseStats.
	cm := statblock.CalculateModifier
	if out.Modifiers.STR != cm(out.FinalStats.STR) {
		t.Errorf("STR modifier: got %d, want %d (final=%d)", out.Modifiers.STR, cm(out.FinalStats.STR), out.FinalStats.STR)
	}
	if out.Modifiers.DEX != cm(out.FinalStats.DEX) {
		t.Errorf("DEX modifier: got %d, want %d (final=%d)", out.Modifiers.DEX, cm(out.FinalStats.DEX), out.FinalStats.DEX)
	}
	if out.Modifiers.CON != cm(out.FinalStats.CON) {
		t.Errorf("CON modifier: got %d, want %d (final=%d)", out.Modifiers.CON, cm(out.FinalStats.CON), out.FinalStats.CON)
	}
	if out.Modifiers.INT != cm(out.FinalStats.INT) {
		t.Errorf("INT modifier: got %d, want %d (final=%d)", out.Modifiers.INT, cm(out.FinalStats.INT), out.FinalStats.INT)
	}
	if out.Modifiers.WIS != cm(out.FinalStats.WIS) {
		t.Errorf("WIS modifier: got %d, want %d (final=%d)", out.Modifiers.WIS, cm(out.FinalStats.WIS), out.FinalStats.WIS)
	}
	if out.Modifiers.CHA != cm(out.FinalStats.CHA) {
		t.Errorf("CHA modifier: got %d, want %d (final=%d)", out.Modifiers.CHA, cm(out.FinalStats.CHA), out.FinalStats.CHA)
	}
}

// ─── 3. FinalStats == BaseStats until background ASI is wired (5.5e stub) ────
// TODO(5.5e): replace this test with background ASI assertions once
// background ASI resolution is implemented in the pipeline.

func TestGenerate_FinalStatsEqualBaseStatsUnderStub(t *testing.T) {
	seed := int64(100)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassBarbarian),
		Species: ptrSpecies(domain.SpeciesHalfOrc),
		Seed:    ptrSeed(seed),
	}

	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// Under the 5.5e stub, no bonuses are applied yet — FinalStats must equal BaseStats.
	if out.FinalStats != out.BaseStats {
		t.Errorf("FinalStats %+v != BaseStats %+v — expected equality under bonus stub", out.FinalStats, out.BaseStats)
	}
}

// ─── 4. HP formula ───────────────────────────────────────────────────────────

func TestGenerate_HPFormula(t *testing.T) {
	hitDieByClass := map[domain.Class]int{
		domain.ClassBarbarian: 12,
		domain.ClassBard:       8,
		domain.ClassCleric:     8,
		domain.ClassDruid:      8,
		domain.ClassFighter:   10,
		domain.ClassMonk:       8,
		domain.ClassPaladin:   10,
		domain.ClassRanger:    10,
		domain.ClassRogue:      8,
		domain.ClassSorcerer:   6,
		domain.ClassWarlock:    8,
		domain.ClassWizard:     6,
		domain.ClassArtificer:  8,
	}

	seed := int64(55)
	for class, hitDie := range hitDieByClass {
		t.Run(string(class), func(t *testing.T) {
			in := statblock.Input{
				Class:   ptrClass(class),
				Species: ptrSpecies(domain.SpeciesHuman),
				Level:   1,
				Seed:    ptrSeed(seed),
			}
			out, err := statblock.Generate(in)
			if err != nil {
				t.Fatalf("Generate: %v", err)
			}
			conMod := out.Modifiers.CON
			wantHP := max(hitDie+conMod, 1)
			if out.Derived.HP != wantHP {
				t.Errorf("HP: got %d, want %d (hitDie=%d conMod=%d)", out.Derived.HP, wantHP, hitDie, conMod)
			}
		})
	}
}

// ─── 5. AC heavy armor — DEX NOT added ───────────────────────────────────────

func TestGenerate_ACHeavyArmor(t *testing.T) {
	seed := int64(99)
	// Fighter uses Chain Mail (heavy, baseAC=16)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassFighter),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out.Derived.AC != 16 {
		t.Errorf("Fighter AC (heavy): got %d, want 16", out.Derived.AC)
	}
}

// ─── 6. AC medium armor — DEX capped at +2 ───────────────────────────────────

func TestGenerate_ACMediumArmor(t *testing.T) {
	seed := int64(1)
	// Cleric uses chain shirt (medium, baseAC=13, maxDex=2)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassCleric),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	dexBonus := min(out.Modifiers.DEX, 2)
	wantAC := 13 + dexBonus
	if out.Derived.AC != wantAC {
		t.Errorf("Cleric AC (medium): got %d, want %d (dexMod=%d)", out.Derived.AC, wantAC, out.Modifiers.DEX)
	}
}

// ─── 7. AC light armor — full DEX ────────────────────────────────────────────

func TestGenerate_ACLightArmor(t *testing.T) {
	seed := int64(2)
	// Rogue uses leather (light, baseAC=11)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassRogue),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	wantAC := 11 + out.Modifiers.DEX
	if out.Derived.AC != wantAC {
		t.Errorf("Rogue AC (light): got %d, want %d (dexMod=%d)", out.Derived.AC, wantAC, out.Modifiers.DEX)
	}
}

// ─── 8. AC unarmored Barbarian — 10 + DEX + CON ─────────────────────────────

func TestGenerate_ACUnarmoredBarbarian(t *testing.T) {
	seed := int64(3)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassBarbarian),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	wantAC := 10 + out.Modifiers.DEX + out.Modifiers.CON
	if out.Derived.AC != wantAC {
		t.Errorf("Barbarian AC (unarmored): got %d, want %d (dex=%d con=%d)", out.Derived.AC, wantAC, out.Modifiers.DEX, out.Modifiers.CON)
	}
}

// ─── 9. AC unarmored Monk — 10 + DEX + WIS ──────────────────────────────────

func TestGenerate_ACUnarmoredMonk(t *testing.T) {
	seed := int64(4)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassMonk),
		Species: ptrSpecies(domain.SpeciesHuman),
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	wantAC := 10 + out.Modifiers.DEX + out.Modifiers.WIS
	if out.Derived.AC != wantAC {
		t.Errorf("Monk AC (unarmored): got %d, want %d (dex=%d wis=%d)", out.Derived.AC, wantAC, out.Modifiers.DEX, out.Modifiers.WIS)
	}
}

// ─── 10. All 13 classes generate without error ───────────────────────────────

func TestGenerate_AllClasses(t *testing.T) {
	classes := []domain.Class{
		domain.ClassBarbarian,
		domain.ClassBard,
		domain.ClassCleric,
		domain.ClassDruid,
		domain.ClassFighter,
		domain.ClassMonk,
		domain.ClassPaladin,
		domain.ClassRanger,
		domain.ClassRogue,
		domain.ClassSorcerer,
		domain.ClassWarlock,
		domain.ClassWizard,
		domain.ClassArtificer,
	}

	seed := int64(77)
	for _, class := range classes {
		t.Run(string(class), func(t *testing.T) {
			in := statblock.Input{
				Class:   ptrClass(class),
				Species: ptrSpecies(domain.SpeciesHuman),
				Seed:    ptrSeed(seed),
			}
			out, err := statblock.Generate(in)
			if err != nil {
				t.Fatalf("Generate: %v", err)
			}
			if out.Class != class {
				t.Errorf("Class: got %q, want %q", out.Class, class)
			}
		})
	}
}

// ─── 11. All 9 species generate without error ────────────────────────────────

func TestGenerate_AllSpecies(t *testing.T) {
	species := []domain.Species{
		domain.SpeciesHuman,
		domain.SpeciesElf,
		domain.SpeciesDwarf,
		domain.SpeciesHalfling,
		domain.SpeciesGnome,
		domain.SpeciesHalfElf,
		domain.SpeciesHalfOrc,
		domain.SpeciesTiefling,
		domain.SpeciesDragonborn,
	}

	seed := int64(88)
	for _, sp := range species {
		t.Run(string(sp), func(t *testing.T) {
			in := statblock.Input{
				Class:   ptrClass(domain.ClassFighter),
				Species: ptrSpecies(sp),
				Seed:    ptrSeed(seed),
			}
			out, err := statblock.Generate(in)
			if err != nil {
				t.Fatalf("Generate: %v", err)
			}
			if out.Species != sp {
				t.Errorf("Species: got %q, want %q", out.Species, sp)
			}
		})
	}
}

// ─── 12. Level defaults to 1 when zero ───────────────────────────────────────

func TestGenerate_LevelDefaultsToOne(t *testing.T) {
	seed := int64(5)
	in := statblock.Input{
		Class:   ptrClass(domain.ClassFighter),
		Species: ptrSpecies(domain.SpeciesHuman),
		Level:   0, // explicitly zero = use default
		Seed:    ptrSeed(seed),
	}
	out, err := statblock.Generate(in)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if out.Level != 1 {
		t.Errorf("Level: got %d, want 1", out.Level)
	}
}

// ─── 13. Nil class picks randomly — always a valid domain.Class ──────────────

func TestGenerate_NilClassPicksRandom(t *testing.T) {
	validClasses := map[domain.Class]bool{
		domain.ClassBarbarian: true,
		domain.ClassBard:      true,
		domain.ClassCleric:    true,
		domain.ClassDruid:     true,
		domain.ClassFighter:   true,
		domain.ClassMonk:      true,
		domain.ClassPaladin:   true,
		domain.ClassRanger:    true,
		domain.ClassRogue:     true,
		domain.ClassSorcerer:  true,
		domain.ClassWarlock:   true,
		domain.ClassWizard:    true,
		domain.ClassArtificer: true,
	}

	for i := range 20 {
		seed := int64(i)
		in := statblock.Input{
			Class:   nil, // no class — should pick randomly
			Species: ptrSpecies(domain.SpeciesHuman),
			Seed:    ptrSeed(seed),
		}
		out, err := statblock.Generate(in)
		if err != nil {
			t.Fatalf("seed %d: Generate: %v", seed, err)
		}
		if !validClasses[out.Class] {
			t.Errorf("seed %d: got invalid class %q", seed, out.Class)
		}
	}
}

// ─── 14. BaseStats always cost exactly 27 in D&D 5e point buy ────────────────

func TestGenerate_BaseStatsPointBuyCost27(t *testing.T) {
	classes := []domain.Class{
		domain.ClassBarbarian, domain.ClassBard, domain.ClassCleric,
		domain.ClassDruid, domain.ClassFighter, domain.ClassMonk,
		domain.ClassPaladin, domain.ClassRanger, domain.ClassRogue,
		domain.ClassSorcerer, domain.ClassWarlock, domain.ClassWizard,
		domain.ClassArtificer,
	}

	for _, class := range classes {
		for i := range 20 {
			seed := int64(i)
			in := statblock.Input{
				Class:   ptrClass(class),
				Species: ptrSpecies(domain.SpeciesHuman),
				Seed:    ptrSeed(seed),
			}
			out, err := statblock.Generate(in)
			if err != nil {
				t.Fatalf("%s seed %d: Generate: %v", class, seed, err)
			}
			cost := statblock.PointBuyCost(out.BaseStats)
			if cost != 27 {
				t.Errorf("%s seed %d: BaseStats point buy cost = %d, want 27 — stats: %+v",
					class, seed, cost, out.BaseStats)
			}
		}
	}
}

// ─── 15. CalculateModifier table test ────────────────────────────────────────

func TestCalculateModifier(t *testing.T) {
	tests := []struct {
		score int
		want  int
	}{
		{10, 0},
		{12, 1},
		{8, -1},
		{20, 5},
		{1, -5},
		{15, 2},
		{11, 0},
		{13, 1},
		{9, -1},
		{18, 4},
		{3, -4},
		{6, -2},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := statblock.CalculateModifier(tt.score)
			if got != tt.want {
				t.Errorf("CalculateModifier(%d) = %d, want %d", tt.score, got, tt.want)
			}
		})
	}
}
