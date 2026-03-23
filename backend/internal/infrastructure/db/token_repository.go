package db

import (
	"context"
	"database/sql"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.TokenRepository = (*TokenRepository)(nil)

// TokenRepository implements ports.TokenRepository using SQLite.
type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// Save inserts a new magic link token record.
func (r *TokenRepository) Save(ctx context.Context, t *domain.MagicLinkToken) error {
	const q = `
		INSERT INTO magic_link_tokens (id, hashed_token, email, expires_at, used_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	var usedAt interface{}
	if t.UsedAt != nil {
		usedAt = t.UsedAt.UTC().Format(time.RFC3339Nano)
	}

	_, err := r.db.ExecContext(ctx, q,
		t.ID,
		t.HashedToken,
		t.Email,
		t.ExpiresAt.UTC().Format(time.RFC3339Nano),
		usedAt,
		t.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

// FindByHashedToken returns the token matching the given hash, or nil if not found.
func (r *TokenRepository) FindByHashedToken(ctx context.Context, hashed string) (*domain.MagicLinkToken, error) {
	const q = `
		SELECT id, hashed_token, email, expires_at, used_at, created_at
		FROM magic_link_tokens
		WHERE hashed_token = ?
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, hashed)
	return scanToken(row)
}

// MarkUsed sets used_at for the token with the given ID.
func (r *TokenRepository) MarkUsed(ctx context.Context, id string, usedAt time.Time) error {
	const q = `UPDATE magic_link_tokens SET used_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, usedAt.UTC().Format(time.RFC3339Nano), id)
	return err
}

func scanToken(row *sql.Row) (*domain.MagicLinkToken, error) {
	var t domain.MagicLinkToken
	var expiresAt, createdAt string
	var usedAt sql.NullString

	err := row.Scan(&t.ID, &t.HashedToken, &t.Email, &expiresAt, &usedAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t.ExpiresAt, err = time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		return nil, err
	}
	t.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}

	if usedAt.Valid {
		parsed, parseErr := time.Parse(time.RFC3339Nano, usedAt.String)
		if parseErr != nil {
			return nil, parseErr
		}
		t.UsedAt = &parsed
	}

	return &t, nil
}
