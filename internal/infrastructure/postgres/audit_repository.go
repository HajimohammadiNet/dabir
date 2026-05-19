package postgres

import (
	"context"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/audit"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log *audit.AuditLog) error {
	const query = `
		INSERT INTO audit_logs (
			actor_user_id,
			action,
			entity_type,
			entity_id,
			old_value,
			new_value,
			ip_address,
			user_agent
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		log.ActorUserID,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.OldValue,
		log.NewValue,
		log.IPAddress,
		log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (r *AuditRepository) List(ctx context.Context, filter audit.ListFilter) ([]audit.AuditLog, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	offset := (filter.Page - 1) * filter.PageSize

	where := "WHERE 1=1"
	args := []interface{}{}
	argPos := 1

	if filter.Action != "" {
		where += fmt.Sprintf(" AND action = $%d", argPos)
		args = append(args, filter.Action)
		argPos++
	}

	if filter.EntityType != "" {
		where += fmt.Sprintf(" AND entity_type = $%d", argPos)
		args = append(args, filter.EntityType)
		argPos++
	}

	if filter.ActorUserID != "" {
		where += fmt.Sprintf(" AND actor_user_id = $%d", argPos)
		args = append(args, filter.ActorUserID)
		argPos++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM audit_logs
		%s
	`, where)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			actor_user_id,
			action,
			entity_type,
			entity_id,
			old_value,
			new_value,
			ip_address::TEXT,
			user_agent,
			created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	items := make([]audit.AuditLog, 0)

	for rows.Next() {
		var item audit.AuditLog

		if err := rows.Scan(
			&item.ID,
			&item.ActorUserID,
			&item.Action,
			&item.EntityType,
			&item.EntityID,
			&item.OldValue,
			&item.NewValue,
			&item.IPAddress,
			&item.UserAgent,
			&item.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate audit logs: %w", err)
	}

	return items, total, nil
}
