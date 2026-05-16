package settings

import "time"

type Setting struct {
	Key       string
	Value     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	KeyOrganizationName = "organization_name"
	KeyLetterConfig     = "letter_config"
)
