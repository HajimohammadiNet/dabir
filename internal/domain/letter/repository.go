package letter

import (
	"context"
	"time"
)

type Repository interface {
	NextNumber(ctx context.Context) (int64, error)

	Create(ctx context.Context, l *Letter) error
	FindByID(ctx context.Context, id string) (*Letter, error)
	List(ctx context.Context, filter ListFilter) ([]Letter, int, error)
	Update(ctx context.Context, l *Letter) error
	SoftDelete(ctx context.Context, id string, deletedBy string) error
}

type ListFilter struct {
	Page     int
	PageSize int

	Search        string
	RegistrarName string
	Sender        string
	Receiver      string

	FromDate *time.Time
	ToDate   *time.Time

	IncludeDeleted bool
}
