package graph

import (
	"forge-rpg/internal/usecase/auth"
	"forge-rpg/internal/usecase/character"
)

// Resolver is the root GraphQL resolver.
// It holds references to all use-case services.
type Resolver struct {
	auth    *auth.Service
	manager *character.Service
}

// NewResolver constructs the root resolver with all required dependencies.
func NewResolver(
	authSvc *auth.Service,
	manager *character.Service,
) *Resolver {
	return &Resolver{
		auth:    authSvc,
		manager: manager,
	}
}
