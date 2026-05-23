package letters

import (
	"fmt"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	domainsettings "github.com/hajimohammadinet/dabir/internal/domain/settings"
	"github.com/hajimohammadinet/dabir/internal/shared/dateutil"
)

type LetterDTO struct {
	ID                    string `json:"id"`
	LetterNumber          int64  `json:"letter_number"`
	FormattedLetterNumber string `json:"formatted_letter_number"`

	Title            string `json:"title"`
	LetterDate       string `json:"letter_date"`
	LetterDateJalali string `json:"letter_date_jalali"`

	RegistrarName string `json:"registrar_name"`
	Sender        string `json:"sender"`
	Receiver      string `json:"receiver"`

	Description *string `json:"description,omitempty"`

	CreatedBy string  `json:"created_by"`
	UpdatedBy *string `json:"updated_by,omitempty"`
	DeletedBy *string `json:"deleted_by,omitempty"`

	IsDeleted bool       `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type NumberingMode string

const (
	NumberingModeFixedPrefix  NumberingMode = "fixed_prefix"
	NumberingModeJalaliYearly NumberingMode = "jalali_yearly"
)

type LetterNumberConfig struct {
	Mode NumberingMode

	Prefix  string
	Padding int

	YearlyPrefixDigits  int
	YearlySerialPadding int
	YearlySeparator     string
	YearSource          domainsettings.YearSource
}

func ToLetterDTO(l letter.Letter, cfg LetterNumberConfig) LetterDTO {
	return LetterDTO{
		ID:                    l.ID,
		LetterNumber:          l.LetterNumber,
		FormattedLetterNumber: FormatLetterNumber(l.LetterNumber, cfg),

		Title:            l.Title,
		LetterDate:       l.LetterDate.Format("2006-01-02"),
		LetterDateJalali: dateutil.ToJalaliString(l.LetterDate),
		RegistrarName:    l.RegistrarName,
		Sender:           l.Sender,
		Receiver:         l.Receiver,
		Description:      l.Description,

		CreatedBy: l.CreatedBy,
		UpdatedBy: l.UpdatedBy,
		DeletedBy: l.DeletedBy,

		IsDeleted: l.IsDeleted,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		DeletedAt: l.DeletedAt,
	}
}

func FormatLetterNumber(number int64, cfg LetterNumberConfig) string {
	padding := cfg.Padding
	if padding <= 0 {
		padding = 6
	}

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "DABIR"
	}

	return fmt.Sprintf("%s-%0*d", prefix, padding, number)
}
