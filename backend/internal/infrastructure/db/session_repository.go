package db

import (
	"context"
	"database/sql"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.SessionRepository = (*SessionRepository)(nil)

// SessionRepository implements ports.SessionRepository using SQLite.
type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Save inserts or replaces a session record.
func (r *SessionRepository) Save(ctx context.Context, s *domain.Session) error {
	const q = `INSERT OR REPLACE INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q,
		s.ID,
		s.UserID,
		s.ExpiresAt.UTC().Format(time.RFC3339Nano),
		s.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

// FindByID returns the session with the given ID, or nil if not found.
func (r *SessionRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	const q = `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, id)

	var s domain.Session
	var expiresAt, createdAt string
	err := row.Scan(&s.ID, &s.UserID, &expiresAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.ExpiresAt, err = time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		return nil, err
	}
	s.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Delete removes the session with the given ID.
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM sessions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
