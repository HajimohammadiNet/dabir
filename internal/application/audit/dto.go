package audit

import (
	"encoding/json"
	"time"

	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type AuditLogDTO struct {
	ID          string  `json:"id"`
	ActorUserID *string `json:"actor_user_id,omitempty"`

	Action     string  `json:"action"`
	EntityType string  `json:"entity_type"`
	EntityID   *string `json:"entity_id,omitempty"`

	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`

	IPAddress *string `json:"ip_address,omitempty"`
	UserAgent *string `json:"user_agent,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func ToAuditLogDTO(log domainaudit.AuditLog) AuditLogDTO {
	return AuditLogDTO{
		ID:          log.ID,
		ActorUserID: log.ActorUserID,
		Action:      log.Action,
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		OldValue:    decodeJSON(log.OldValue),
		NewValue:    decodeJSON(log.NewValue),
		IPAddress:   log.IPAddress,
		UserAgent:   log.UserAgent,
		CreatedAt:   log.CreatedAt,
	}
}

func decodeJSON(raw []byte) interface{} {
	if len(raw) == 0 {
		return nil
	}

	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return string(raw)
	}

	return value
}
