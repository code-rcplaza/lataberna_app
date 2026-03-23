package ports

import (
	"context"
	"time"

	"forge-rpg/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
}

type SessionRepository interface {
	Save(ctx context.Context, s *domain.Session) error
	FindByID(ctx context.Context, id string) (*domain.Session, error)
	Delete(ctx context.Context, id string) error
}

type TokenRepository interface {
	Save(ctx context.Context, t *domain.MagicLinkToken) error
	FindByHashedToken(ctx context.Context, hashed string) (*domain.MagicLinkToken, error)
	MarkUsed(ctx context.Context, id string, usedAt time.Time) error
}

// Mailer sends transactional emails.
type Mailer interface {
	SendMagicLink(ctx context.Context, email, magicLink string) error
}
