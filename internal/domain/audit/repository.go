package audit

import "context"

type Repository interface {
	Create(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, filter ListFilter) ([]AuditLog, int, error)
}

type ListFilter struct {
	Page     int
	PageSize int

	Action      string
	EntityType  string
	ActorUserID string
}
