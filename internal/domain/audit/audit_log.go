package audit

import "time"

type AuditLog struct {
	ID          string
	ActorUserID *string

	Action     string
	EntityType string
	EntityID   *string

	OldValue []byte
	NewValue []byte

	IPAddress *string
	UserAgent *string

	CreatedAt time.Time
}

const (
	ActionSetupInitialized = "setup.initialized"

	ActionAuthLoginSuccess = "auth.login_success"
	ActionAuthLoginFailed  = "auth.login_failed"

	ActionUserCreated     = "user.created"
	ActionUserUpdated     = "user.updated"
	ActionUserActivated   = "user.activated"
	ActionUserDeactivated = "user.deactivated"

	ActionLetterCreated = "letter.created"
	ActionLetterUpdated = "letter.updated"
	ActionLetterDeleted = "letter.deleted"

	ActionLettersImportPreviewed = "letters.import_previewed"
	ActionLettersImportCommitted = "letters.import_committed"
)
