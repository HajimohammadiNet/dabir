package importjob

import "context"

type Repository interface {
	Create(ctx context.Context, job *ImportJob) error
	FindByID(ctx context.Context, id string) (*ImportJob, error)
	MarkCommitted(ctx context.Context, id string, committedBy string) error
}
