package domain

// ArmorCategory defines the type of armor.
type ArmorCategory string

const (
	ArmorNone   ArmorCategory = "none"
	ArmorLight  ArmorCategory = "light"
	ArmorMedium ArmorCategory = "medium"
	ArmorHeavy  ArmorCategory = "heavy"
	ArmorShield ArmorCategory = "shield"
)

// ArmorType is a resource that describes a piece of armor.
// It is NOT embedded in class logic — classes reference it by name.
type ArmorType struct {
	Name           string
	Category       ArmorCategory
	BaseAC         int
	MaxDex         *int // nil = unlimited, 2 = medium armor, 0 = heavy armor
	StrRequirement *int
}
