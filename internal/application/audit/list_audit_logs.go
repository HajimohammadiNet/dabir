package audit

import (
	"context"
	"strings"

	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type ListAuditLogsUseCase struct {
	auditRepo domainaudit.Repository
}

type ListAuditLogsInput struct {
	Page     int
	PageSize int

	Action      string
	EntityType  string
	ActorUserID string
}

type ListAuditLogsOutput struct {
	Items      []AuditLogDTO `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func NewListAuditLogsUseCase(auditRepo domainaudit.Repository) *ListAuditLogsUseCase {
	return &ListAuditLogsUseCase{
		auditRepo: auditRepo,
	}
}

func (uc *ListAuditLogsUseCase) Execute(ctx context.Context, input ListAuditLogsInput) (*ListAuditLogsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	input.Action = strings.TrimSpace(input.Action)
	input.EntityType = strings.TrimSpace(input.EntityType)
	input.ActorUserID = strings.TrimSpace(input.ActorUserID)

	items, total, err := uc.auditRepo.List(ctx, domainaudit.ListFilter{
		Page:        input.Page,
		PageSize:    input.PageSize,
		Action:      input.Action,
		EntityType:  input.EntityType,
		ActorUserID: input.ActorUserID,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]AuditLogDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, ToAuditLogDTO(item))
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + input.PageSize - 1) / input.PageSize
	}

	return &ListAuditLogsOutput{
		Items:      dtos,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}
