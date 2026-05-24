package letters

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	domainsettings "github.com/hajimohammadinet/dabir/internal/domain/settings"
	"github.com/hajimohammadinet/dabir/internal/shared/dateutil"
)

type LetterDTO struct {
	ID                    string `json:"id"`
	LetterNumber          int64  `json:"letter_number"`
	FormattedLetterNumber string `json:"formatted_letter_number"`

	DisplayLetterNumber *string `json:"display_letter_number,omitempty"`

	LetterYear       *int    `json:"letter_year,omitempty"`
	LetterYearSuffix *string `json:"letter_year_suffix,omitempty"`
	LetterSerial     *int64  `json:"letter_serial,omitempty"`

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
	NumberingModeManual       NumberingMode = "manual"
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
		FormattedLetterNumber: FormatLetterNumber(l, cfg),
		DisplayLetterNumber:   l.DisplayLetterNumber,

		LetterYear:       l.LetterYear,
		LetterYearSuffix: l.LetterYearSuffix,
		LetterSerial:     l.LetterSerial,

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

func FormatLetterNumber(l letter.Letter, cfg LetterNumberConfig) string {

	if l.DisplayLetterNumber != nil && strings.TrimSpace(*l.DisplayLetterNumber) != "" {
		return strings.TrimSpace(*l.DisplayLetterNumber)
	}

	if cfg.Mode == NumberingModeJalaliYearly && l.LetterSerial != nil {
		separator := cfg.YearlySeparator
		if separator == "" {
			separator = "-"
		}

		serialPadding := cfg.YearlySerialPadding
		if serialPadding <= 0 {
			serialPadding = 4
		}

		yearSuffix := ""
		if l.LetterYearSuffix != nil && *l.LetterYearSuffix != "" {
			yearSuffix = *l.LetterYearSuffix
		} else if l.LetterYear != nil {
			yearSuffix = BuildJalaliYearSuffix(*l.LetterYear, cfg.YearlyPrefixDigits)
		}

		if yearSuffix != "" {
			return fmt.Sprintf("%s%s%0*d", yearSuffix, separator, serialPadding, *l.LetterSerial)
		}
	}

	padding := cfg.Padding
	if padding <= 0 {
		padding = 6
	}

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "DABIR"
	}

	return fmt.Sprintf("%s-%0*d", prefix, padding, l.LetterNumber)
}

func BuildJalaliYearSuffix(jalaliYear int, digits int) string {
	if digits <= 0 {
		digits = 3
	}

	divisor := int(math.Pow10(digits))
	suffix := jalaliYear % divisor

	return fmt.Sprintf("%0*d", digits, suffix)
}
