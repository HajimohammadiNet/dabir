package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/attachment"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachmentRepository struct {
	db *pgxpool.Pool
}

func NewAttachmentRepository(db *pgxpool.Pool) *AttachmentRepository {
	return &AttachmentRepository{
		db: db,
	}
}

func (r *AttachmentRepository) Create(ctx context.Context, item *attachment.Attachment) error {
	const query = `
		INSERT INTO letter_attachments (
			letter_id,
			file_name,
			object_key,
			content_type,
			size_bytes,
			uploaded_by,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, false)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		item.LetterID,
		item.FileName,
		item.ObjectKey,
		item.ContentType,
		item.SizeBytes,
		item.UploadedBy,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create attachment: %w", err)
	}

	return nil
}

func (r *AttachmentRepository) ListByLetterID(ctx context.Context, letterID string) ([]attachment.Attachment, error) {
	const query = `
		SELECT
			id,
			letter_id,
			file_name,
			object_key,
			content_type,
			size_bytes,
			uploaded_by,
			is_deleted,
			deleted_by,
			deleted_at,
			created_at,
			updated_at
		FROM letter_attachments
		WHERE letter_id = $1
		  AND is_deleted = false
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query, letterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	defer rows.Close()

	items := make([]attachment.Attachment, 0)

	for rows.Next() {
		var item attachment.Attachment

		if err := rows.Scan(
			&item.ID,
			&item.LetterID,
			&item.FileName,
			&item.ObjectKey,
			&item.ContentType,
			&item.SizeBytes,
			&item.UploadedBy,
			&item.IsDeleted,
			&item.DeletedBy,
			&item.DeletedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate attachments: %w", err)
	}

	return items, nil
}

func (r *AttachmentRepository) FindByID(ctx context.Context, id string) (*attachment.Attachment, error) {
	const query = `
		SELECT
			id,
			letter_id,
			file_name,
			object_key,
			content_type,
			size_bytes,
			uploaded_by,
			is_deleted,
			deleted_by,
			deleted_at,
			created_at,
			updated_at
		FROM letter_attachments
		WHERE id = $1
		LIMIT 1
	`

	var item attachment.Attachment

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.LetterID,
		&item.FileName,
		&item.ObjectKey,
		&item.ContentType,
		&item.SizeBytes,
		&item.UploadedBy,
		&item.IsDeleted,
		&item.DeletedBy,
		&item.DeletedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find attachment: %w", err)
	}

	return &item, nil
}

func (r *AttachmentRepository) SoftDelete(ctx context.Context, id string, deletedBy string) error {
	const query = `
		UPDATE letter_attachments
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
		return fmt.Errorf("failed to soft delete attachment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
