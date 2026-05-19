package importjob

import "time"

type Type string
type Status string

const (
	TypeLetters Type = "letters"

	StatusPreviewed Status = "previewed"
	StatusCommitted Status = "committed"
	StatusFailed    Status = "failed"
)

type ImportJob struct {
	ID string

	Type   Type
	Status Status

	FileName string

	TotalRows   int
	ValidRows   int
	InvalidRows int

	MaxLetterNumber *int64

	DetectedColumns []byte
	PreviewData     []byte
	Errors          []byte

	CreatedBy   string
	CommittedBy *string

	CreatedAt   time.Time
	CommittedAt *time.Time
}
