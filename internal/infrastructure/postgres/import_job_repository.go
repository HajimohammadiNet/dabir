package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/importjob"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImportJobRepository struct {
	db *pgxpool.Pool
}

func NewImportJobRepository(db *pgxpool.Pool) *ImportJobRepository {
	return &ImportJobRepository{
		db: db,
	}
}

func (r *ImportJobRepository) Create(ctx context.Context, job *importjob.ImportJob) error {
	const query = `
		INSERT INTO import_jobs (
			type,
			status,
			file_name,
			total_rows,
			valid_rows,
			invalid_rows,
			max_letter_number,
			detected_columns,
			preview_data,
			errors,
			created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		string(job.Type),
		string(job.Status),
		job.FileName,
		job.TotalRows,
		job.ValidRows,
		job.InvalidRows,
		job.MaxLetterNumber,
		job.DetectedColumns,
		job.PreviewData,
		job.Errors,
		job.CreatedBy,
	).Scan(&job.ID, &job.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create import job: %w", err)
	}

	return nil
}

func (r *ImportJobRepository) FindByID(ctx context.Context, id string) (*importjob.ImportJob, error) {
	const query = `
		SELECT
			id,
			type,
			status,
			file_name,
			total_rows,
			valid_rows,
			invalid_rows,
			max_letter_number,
			detected_columns,
			preview_data,
			errors,
			created_by,
			committed_by,
			created_at,
			committed_at
		FROM import_jobs
		WHERE id = $1
		LIMIT 1
	`

	job := &importjob.ImportJob{}
	var jobType string
	var status string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&jobType,
		&status,
		&job.FileName,
		&job.TotalRows,
		&job.ValidRows,
		&job.InvalidRows,
		&job.MaxLetterNumber,
		&job.DetectedColumns,
		&job.PreviewData,
		&job.Errors,
		&job.CreatedBy,
		&job.CommittedBy,
		&job.CreatedAt,
		&job.CommittedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find import job: %w", err)
	}

	job.Type = importjob.Type(jobType)
	job.Status = importjob.Status(status)

	return job, nil
}

func (r *ImportJobRepository) MarkCommitted(ctx context.Context, id string, committedBy string) error {
	const query = `
		UPDATE import_jobs
		SET
			status = $1,
			committed_by = $2,
			committed_at = NOW()
		WHERE id = $3
		  AND status = $4
	`

	result, err := r.db.Exec(
		ctx,
		query,
		string(importjob.StatusCommitted),
		committedBy,
		id,
		string(importjob.StatusPreviewed),
	)
	if err != nil {
		return fmt.Errorf("failed to mark import job committed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
