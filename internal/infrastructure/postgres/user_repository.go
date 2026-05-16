package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hajimohammadinet/dabir/internal/domain/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CountSuperUsers(ctx context.Context) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM users
		WHERE role = 'superuser'
		  AND is_active = true
	`

	var count int

	if err := r.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count superusers: %w", err)
	}

	return count, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	const query = `
		INSERT INTO users (
			username,
			full_name,
			password_hash,
			role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		u.Username,
		u.FullName,
		u.PasswordHash,
		string(u.Role),
		u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	const query = `
		SELECT
			id,
			username,
			full_name,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	u := &user.User{}

	var role string

	err := r.db.QueryRow(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.PasswordHash,
		&role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	u.Role = user.Role(role)

	return u, nil
}
