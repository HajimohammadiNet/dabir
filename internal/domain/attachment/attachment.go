package attachment

import (
	"context"
	"time"
)

type Attachment struct {
	ID string

	LetterID string

	FileName    string
	ObjectKey   string
	ContentType string
	SizeBytes   int64

	UploadedBy string

	IsDeleted bool
	DeletedBy *string
	DeletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, attachment *Attachment) error
	ListByLetterID(ctx context.Context, letterID string) ([]Attachment, error)
	FindByID(ctx context.Context, id string) (*Attachment, error)
	SoftDelete(ctx context.Context, id string, deletedBy string) error
}
