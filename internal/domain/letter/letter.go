package letter

import "time"

type Direction string

const (
	DirectionIncoming Direction = "incoming"
	DirectionOutgoing Direction = "outgoing"
)

func (d Direction) IsValid() bool {
	return d == DirectionIncoming || d == DirectionOutgoing
}

type Letter struct {
	ID           string
	Direction    Direction
	LetterNumber int64

	DisplayLetterNumber *string

	LetterYear       *int
	LetterYearSuffix *string
	LetterSerial     *int64

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
