package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LetterRepository struct {
	db *pgxpool.Pool
}

func NewLetterRepository(db *pgxpool.Pool) *LetterRepository {
	return &LetterRepository{
		db: db,
	}
}

func (r *LetterRepository) NextNumber(ctx context.Context) (int64, error) {
	const query = `SELECT nextval('letter_number_seq')`

	var number int64
	if err := r.db.QueryRow(ctx, query).Scan(&number); err != nil {
		return 0, fmt.Errorf("failed to get next letter number: %w", err)
	}

	return number, nil
}

func (r *LetterRepository) Create(ctx context.Context, l *letter.Letter) error {
	const query = `
		INSERT INTO letters (
			letter_number,
			title,
			letter_date,
			registrar_name,
			sender,
			receiver,
			destination,
			description,
			created_by,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $6, $7, $8, false)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		l.LetterNumber,
		l.Title,
		l.LetterDate,
		l.RegistrarName,
		l.Sender,
		l.Receiver,
		l.Description,
		l.CreatedBy,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create letter: %w", err)
	}

	return nil
}

func (r *LetterRepository) FindByID(ctx context.Context, id string) (*letter.Letter, error) {
	const query = `
		SELECT
			id,
			letter_number,
			title,
			letter_date,
			registrar_name,
			sender,
			receiver,
			destination,
			description,
			created_by,
			updated_by,
			deleted_by,
			is_deleted,
			created_at,
			updated_at,
			deleted_at
		FROM letters
		WHERE id = $1
		LIMIT 1
	`

	l := &letter.Letter{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&l.ID,
		&l.LetterNumber,
		&l.Title,
		&l.LetterDate,
		&l.RegistrarName,
		&l.Sender,
		&l.Receiver,
		&l.Description,
		&l.CreatedBy,
		&l.UpdatedBy,
		&l.DeletedBy,
		&l.IsDeleted,
		&l.CreatedAt,
		&l.UpdatedAt,
		&l.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find letter by id: %w", err)
	}

	return l, nil
}

func (r *LetterRepository) List(ctx context.Context, filter letter.ListFilter) ([]letter.Letter, int, error) {
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

	if !filter.IncludeDeleted {
		where += " AND is_deleted = false"
	}

	if filter.Search != "" {
		where += fmt.Sprintf(
			" AND (title ILIKE $%d OR sender ILIKE $%d OR receiver ILIKE $%d OR registrar_name ILIKE $%d OR CAST(letter_number AS TEXT) ILIKE $%d)",
			argPos,
			argPos,
			argPos,
			argPos,
			argPos,
		)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	if filter.RegistrarName != "" {
		where += fmt.Sprintf(" AND registrar_name ILIKE $%d", argPos)
		args = append(args, "%"+filter.RegistrarName+"%")
		argPos++
	}

	if filter.Sender != "" {
		where += fmt.Sprintf(" AND sender ILIKE $%d", argPos)
		args = append(args, "%"+filter.Sender+"%")
		argPos++
	}

	if filter.Receiver != "" {
		where += fmt.Sprintf(" AND receiver ILIKE $%d", argPos)
		args = append(args, "%"+filter.Receiver+"%")
		argPos++
	}

	if filter.FromDate != nil {
		where += fmt.Sprintf(" AND letter_date >= $%d", argPos)
		args = append(args, *filter.FromDate)
		argPos++
	}

	if filter.ToDate != nil {
		where += fmt.Sprintf(" AND letter_date <= $%d", argPos)
		args = append(args, *filter.ToDate)
		argPos++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM letters
		%s
	`, where)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count letters: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			letter_number,
			title,
			letter_date,
			registrar_name,
			sender,
			receiver,
			destination,
			description,
			created_by,
			updated_by,
			deleted_by,
			is_deleted,
			created_at,
			updated_at,
			deleted_at
		FROM letters
		%s
		ORDER BY letter_number DESC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list letters: %w", err)
	}
	defer rows.Close()

	items := make([]letter.Letter, 0)

	for rows.Next() {
		var l letter.Letter

		if err := rows.Scan(
			&l.ID,
			&l.LetterNumber,
			&l.Title,
			&l.LetterDate,
			&l.RegistrarName,
			&l.Sender,
			&l.Receiver,
			&l.Description,
			&l.CreatedBy,
			&l.UpdatedBy,
			&l.DeletedBy,
			&l.IsDeleted,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan letter: %w", err)
		}

		items = append(items, l)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate letters: %w", err)
	}

	return items, total, nil
}

func (r *LetterRepository) Update(ctx context.Context, l *letter.Letter) error {
	const query = `
		UPDATE letters
		SET
			title = $1,
			letter_date = $2,
			sender = $3,
			receiver = $4,
			destination = $4,
			description = $5,
			updated_by = $6,
			updated_at = NOW()
		WHERE id = $7
		AND is_deleted = false
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		l.Title,
		l.LetterDate,
		l.Sender,
		l.Receiver,
		l.Description,
		l.UpdatedBy,
		l.ID,
	).Scan(&l.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}

		return fmt.Errorf("failed to update letter: %w", err)
	}

	return nil
}

func (r *LetterRepository) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	const query = `
		UPDATE letters
		SET
			is_deleted = true,
			deleted_by = $1,
			deleted_at = $2,
			updated_at = $2
		WHERE id = $3
		  AND is_deleted = false
	`

	now := time.Now().UTC()

	result, err := r.db.Exec(ctx, query, deletedBy, now, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete letter: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
