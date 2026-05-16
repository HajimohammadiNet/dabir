package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hajimohammadinet/dabir/internal/domain/settings"
)

type SettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{
		db: db,
	}
}

func (r *SettingsRepository) Set(ctx context.Context, key string, value []byte) error {
	const query = `
		INSERT INTO app_settings (
			key,
			value
		)
		VALUES ($1, $2)
		ON CONFLICT (key)
		DO UPDATE SET
			value = EXCLUDED.value,
			updated_at = NOW()
	`

	if _, err := r.db.Exec(ctx, query, key, value); err != nil {
		return fmt.Errorf("failed to set app setting: %w", err)
	}

	return nil
}

func (r *SettingsRepository) Get(ctx context.Context, key string) (*settings.Setting, error) {
	const query = `
		SELECT
			key,
			value,
			created_at,
			updated_at
		FROM app_settings
		WHERE key = $1
		LIMIT 1
	`

	s := &settings.Setting{}

	err := r.db.QueryRow(ctx, query, key).Scan(
		&s.Key,
		&s.Value,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get app setting: %w", err)
	}

	return s, nil
}
