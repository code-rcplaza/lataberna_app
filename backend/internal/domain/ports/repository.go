package ports

import (
	"context"
	"forge-rpg/internal/domain"
)

// CharacterRepository defines the persistence interface for characters.
// The infrastructure layer implements this; the domain never imports infrastructure.
type CharacterRepository interface {
	Save(ctx context.Context, c *domain.Character) error
	FindByID(ctx context.Context, id string) (*domain.Character, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Character, error)
	Update(ctx context.Context, c *domain.Character) error
	Delete(ctx context.Context, id string) error
}
