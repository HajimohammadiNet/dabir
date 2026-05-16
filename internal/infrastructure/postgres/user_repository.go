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

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
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
		WHERE id = $1
		LIMIT 1
	`

	u := &user.User{}
	var role string

	err := r.db.QueryRow(ctx, query, id).Scan(
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

		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	u.Role = user.Role(role)

	return u, nil
}

func (r *UserRepository) List(ctx context.Context, filter user.ListFilter) ([]user.User, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	where := "WHERE 1=1"
	args := []interface{}{}
	argPos := 1

	if filter.Search != "" {
		where += fmt.Sprintf(" AND (username ILIKE $%d OR full_name ILIKE $%d)", argPos, argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	if filter.Role != nil {
		where += fmt.Sprintf(" AND role = $%d", argPos)
		args = append(args, string(*filter.Role))
		argPos++
	}

	if filter.IsActive != nil {
		where += fmt.Sprintf(" AND is_active = $%d", argPos)
		args = append(args, *filter.IsActive)
		argPos++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM users
		%s
	`, where)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	query := fmt.Sprintf(`
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
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	items := make([]user.User, 0)

	for rows.Next() {
		u := user.User{}
		var role string

		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FullName,
			&u.PasswordHash,
			&role,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		u.Role = user.Role(role)
		items = append(items, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate users: %w", err)
	}

	return items, total, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	const query = `
		UPDATE users
		SET
			full_name = $1,
			role = $2,
			is_active = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	if err := r.db.QueryRow(
		ctx,
		query,
		u.FullName,
		string(u.Role),
		u.IsActive,
		u.ID,
	).Scan(&u.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) SetActive(ctx context.Context, id string, isActive bool) error {
	const query = `
		UPDATE users
		SET
			is_active = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, isActive, id)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
