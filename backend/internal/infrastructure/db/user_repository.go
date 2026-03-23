package db

import (
	"context"
	"database/sql"
	"time"

	"forge-rpg/internal/domain"
	"forge-rpg/internal/domain/ports"
)

// Compile-time interface check.
var _ ports.UserRepository = (*UserRepository)(nil)

// UserRepository implements ports.UserRepository using SQLite.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByEmail returns the user with the given email, or nil if not found.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `SELECT id, email, created_at FROM users WHERE email = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, email)
	return scanUser(row)
}

// FindByID returns the user with the given ID, or nil if not found.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `SELECT id, email, created_at FROM users WHERE id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, id)
	return scanUser(row)
}

// Save inserts or replaces a user record.
func (r *UserRepository) Save(ctx context.Context, u *domain.User) error {
	const q = `INSERT OR REPLACE INTO users (id, email, created_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q, u.ID, u.Email, u.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}

func scanUser(row *sql.Row) (*domain.User, error) {
	var u domain.User
	var createdAt string
	err := row.Scan(&u.ID, &u.Email, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = t
	return &u, nil
}
