package domain

// NarrativeCategory identifies which narrative block this is.
type NarrativeCategory string

const (
	NarrativeBackground NarrativeCategory = "background"
	NarrativeMotivation NarrativeCategory = "motivation"
	NarrativeSecret     NarrativeCategory = "secret"
)

// NarrativeBlock is a single narrative entry (background, motivation, or secret).
// Tags control which class/species combinations can use this block.
// The tag "any" means universally compatible.
type NarrativeBlock struct {
	Category NarrativeCategory
	Content  string
	Tags     []string
}
