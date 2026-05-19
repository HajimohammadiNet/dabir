package audit

import (
	"context"
	"encoding/json"

	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type Logger struct {
	auditRepo domainaudit.Repository
}

type LogInput struct {
	ActorUserID *string

	Action     string
	EntityType string
	EntityID   *string

	OldValue interface{}
	NewValue interface{}

	IPAddress *string
	UserAgent *string
}

func NewLogger(auditRepo domainaudit.Repository) *Logger {
	return &Logger{
		auditRepo: auditRepo,
	}
}

func (l *Logger) Log(ctx context.Context, input LogInput) {
	oldValue := marshalOrNil(input.OldValue)
	newValue := marshalOrNil(input.NewValue)

	_ = l.auditRepo.Create(ctx, &domainaudit.AuditLog{
		ActorUserID: input.ActorUserID,
		Action:      input.Action,
		EntityType:  input.EntityType,
		EntityID:    input.EntityID,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   input.IPAddress,
		UserAgent:   input.UserAgent,
	})
}

func marshalOrNil(value interface{}) []byte {
	if value == nil {
		return nil
	}

	b, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	return b
}
