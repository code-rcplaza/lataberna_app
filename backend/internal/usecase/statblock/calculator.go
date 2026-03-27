package statblock

import (
	"math/rand"

	"forge-rpg/internal/domain"
)

// CalculateModifier returns ⌊(score - 10) / 2⌋ using true floor division.
// Go's integer division truncates toward zero, so we need explicit floor for negatives.
// This is exported so tests can verify it directly.
// It MUST always be called on FinalStats values, never BaseStats.
func CalculateModifier(score int) int {
	n := score - 10
	// Floor division: for negative odd numbers, truncation gives wrong result.
	if n < 0 && n%2 != 0 {
		return (n - 1) / 2
	}
	return n / 2
}

func calculateModifiers(s domain.Stats) domain.Modifiers {
	return domain.Modifiers{
		STR: CalculateModifier(s.STR),
		DEX: CalculateModifier(s.DEX),
		CON: CalculateModifier(s.CON),
		INT: CalculateModifier(s.INT),
		WIS: CalculateModifier(s.WIS),
		CHA: CalculateModifier(s.CHA),
	}
}

// calculateHP computes HP at level 1: full hit die + CON modifier.
// Level parameter is accepted now so future levels can extend this.
func calculateHP(class domain.Class, mods domain.Modifiers, level int) int {
	hitDie := hitDieFor(class)
	// Level 1: maximum hit die value + CON modifier
	hp := hitDie + mods.CON
	if hp < 1 {
		hp = 1 // HP minimum is always 1
	}
	return hp
}

func hitDieFor(class domain.Class) int {
	data, ok := classTable[class]
	if !ok {
		return 8 // safe fallback
	}
	return data.hitDie
}

// calculateAC computes AC based on armor category and modifiers.
func calculateAC(armor domain.ArmorType, mods domain.Modifiers) int {
	switch armor.Category {
	case domain.ArmorHeavy:
		return armor.BaseAC // no DEX modifier

	case domain.ArmorMedium:
		dexBonus := mods.DEX
		if armor.MaxDex != nil && dexBonus > *armor.MaxDex {
			dexBonus = *armor.MaxDex
		}
		return armor.BaseAC + dexBonus

	case domain.ArmorLight:
		return armor.BaseAC + mods.DEX

	case domain.ArmorNone:
		return armor.BaseAC + mods.DEX

	case domain.ArmorCategory("unarmored-barbarian"):
		return 10 + mods.DEX + mods.CON

	case domain.ArmorCategory("unarmored-monk"):
		return 10 + mods.DEX + mods.WIS

	default:
		return armor.BaseAC
	}
}

// pointBuyCostTable maps a stat score to its D&D 5e point buy cost.
// Valid range is 8–15; scores outside this range are not in the table.
var pointBuyCostTable = map[int]int{
	8: 0, 9: 1, 10: 2, 11: 3, 12: 4, 13: 5, 14: 7, 15: 9,
}

// PointBuyCost returns the total D&D 5e point buy cost for a stat block.
// Returns -1 if any stat is outside the valid point buy range [8, 15].
// Exported so tests can assert the 27-point invariant.
func PointBuyCost(s domain.Stats) int {
	total := 0
	for _, v := range []int{s.STR, s.DEX, s.CON, s.INT, s.WIS, s.CHA} {
		c, ok := pointBuyCostTable[v]
		if !ok {
			return -1
		}
		total += c
	}
	return total
}

// buildStatsFromPriority generates base stats using the D&D 5e standard array.
// primary receives 15, secondary receives 14.
// The remaining four values [13, 12, 10, 8] are randomly shuffled into the
// other four stats. Total point buy cost is always exactly 27.
func buildStatsFromPriority(primary, secondary string, rng *rand.Rand) domain.Stats {
	remaining := [4]int{13, 12, 10, 8}
	rng.Shuffle(len(remaining), func(i, j int) {
		remaining[i], remaining[j] = remaining[j], remaining[i]
	})

	vals := map[string]int{
		"STR": 0, "DEX": 0, "CON": 0, "INT": 0, "WIS": 0, "CHA": 0,
	}
	vals[primary] = 15
	vals[secondary] = 14

	ri := 0
	for _, stat := range []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"} {
		if vals[stat] == 0 {
			vals[stat] = remaining[ri]
			ri++
		}
	}

	return domain.Stats{
		STR: vals["STR"],
		DEX: vals["DEX"],
		CON: vals["CON"],
		INT: vals["INT"],
		WIS: vals["WIS"],
		CHA: vals["CHA"],
	}
}
