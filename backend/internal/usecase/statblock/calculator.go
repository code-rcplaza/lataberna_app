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

// applyVariation adds a controlled random variation of ±1 or ±2 to each stat.
// Stats are clamped to [6, 18].
func applyVariation(base domain.Stats, rng *rand.Rand) domain.Stats {
	return domain.Stats{
		STR: clamp(base.STR+variation(rng), 6, 18),
		DEX: clamp(base.DEX+variation(rng), 6, 18),
		CON: clamp(base.CON+variation(rng), 6, 18),
		INT: clamp(base.INT+variation(rng), 6, 18),
		WIS: clamp(base.WIS+variation(rng), 6, 18),
		CHA: clamp(base.CHA+variation(rng), 6, 18),
	}
}

func variation(rng *rand.Rand) int {
	v := []int{-2, -1, -1, 0, 0, 1, 1, 2}
	return v[rng.Intn(len(v))]
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
