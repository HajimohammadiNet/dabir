package letter

import (
	"context"
	"time"
)

type Repository interface {
	NextNumber(ctx context.Context) (int64, error)
	NextNumberForYear(ctx context.Context, jalaliYear int) (int64, error)

	ExistsByDisplayLetterNumber(ctx context.Context, displayNumber string) (bool, error)

	Create(ctx context.Context, l *Letter) error
	FindByID(ctx context.Context, id string) (*Letter, error)
	List(ctx context.Context, filter ListFilter) ([]Letter, int, error)
	Update(ctx context.Context, l *Letter) error
	SoftDelete(ctx context.Context, id string, deletedBy string) error

	BulkCreate(ctx context.Context, letters []Letter) error
	SetSequenceValue(ctx context.Context, value int64) error
	FindExistingNumbers(ctx context.Context, numbers []int64) (map[int64]bool, error)
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

	SortBy    string
	SortOrder string
}
