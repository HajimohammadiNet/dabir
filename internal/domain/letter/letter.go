package letter

import "time"

type Letter struct {
	ID           string
	LetterNumber int64

	Title      string
	LetterDate time.Time

	RegistrarName string
	Sender        string
	Receiver      string

	Description *string

	CreatedBy string
	UpdatedBy *string
	DeletedBy *string

	IsDeleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
