package letters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
	"github.com/hajimohammadinet/dabir/internal/shared/dateutil"
)

var (
	ErrLetterNotFound = errors.New("letter not found")
	ErrInvalidInput   = errors.New("invalid input")
)

type CreateLetterUseCase struct {
	letterRepo     letter.Repository
	configProvider *LetterConfigProvider
}

type CreateLetterInput struct {
	Title      string `json:"title"`
	LetterDate string `json:"letter_date"`

	DisplayLetterNumber *string `json:"display_letter_number"`

	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`

	Description *string `json:"description"`

	ActorUserID   string `json:"-"`
	RegistrarName string `json:"-"`
}

func NewCreateLetterUseCase(
	letterRepo letter.Repository,
	configProvider *LetterConfigProvider,
) *CreateLetterUseCase {
	return &CreateLetterUseCase{
		letterRepo:     letterRepo,
		configProvider: configProvider,
	}
}

func (uc *CreateLetterUseCase) Execute(ctx context.Context, input CreateLetterInput) (*LetterDTO, error) {
	input = normalizeCreateLetterInput(input)

	if err := validateCreateLetterInput(input); err != nil {
		return nil, err
	}

	letterDate, err := dateutil.ParseOfficialDate(input.LetterDate)
	if err != nil {
		return nil, errors.New("letter_date must be in Jalali YYYY/MM/DD format")
	}

	cfg := uc.configProvider.Get(ctx)

	l := &letter.Letter{
		Title:      input.Title,
		LetterDate: letterDate,

		RegistrarName: input.RegistrarName,
		Sender:        input.Sender,
		Receiver:      input.Receiver,

		Description: normalizeOptionalString(input.Description),

		CreatedBy: input.ActorUserID,
		IsDeleted: false,
	}

	if cfg.Mode == NumberingModeManual {
		if input.DisplayLetterNumber == nil || *input.DisplayLetterNumber == "" {
			return nil, errors.New("display_letter_number is required in manual numbering mode")
		}

		exists, err := uc.letterRepo.ExistsByDisplayLetterNumber(ctx, *input.DisplayLetterNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to check display letter number uniqueness: %w", err)
		}

		if exists {
			return nil, errors.New("display_letter_number already exists")
		}

		nextNumber, err := uc.letterRepo.NextNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate internal letter number: %w", err)
		}

		l.LetterNumber = nextNumber
		l.DisplayLetterNumber = input.DisplayLetterNumber
	} else if cfg.Mode == NumberingModeJalaliYearly {
		jalaliYear := resolveJalaliYear(letterDate, cfg)
		yearSuffix := BuildJalaliYearSuffix(jalaliYear, cfg.YearlyPrefixDigits)

		nextSerial, err := uc.letterRepo.NextNumberForYear(ctx, jalaliYear)
		if err != nil {
			return nil, fmt.Errorf("failed to generate yearly letter number: %w", err)
		}

		l.LetterNumber = nextSerial
		l.LetterYear = &jalaliYear
		l.LetterYearSuffix = &yearSuffix
		l.LetterSerial = &nextSerial
	} else {
		nextNumber, err := uc.letterRepo.NextNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate letter number: %w", err)
		}

		l.LetterNumber = nextNumber
	}

	if err := uc.letterRepo.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("failed to create letter: %w", err)
	}

	dto := ToLetterDTO(*l, cfg)

	return &dto, nil
}

func normalizeCreateLetterInput(input CreateLetterInput) CreateLetterInput {
	input.Title = strings.TrimSpace(input.Title)
	input.LetterDate = strings.TrimSpace(input.LetterDate)
	input.Sender = strings.TrimSpace(input.Sender)
	input.Receiver = strings.TrimSpace(input.Receiver)
	input.RegistrarName = strings.TrimSpace(input.RegistrarName)
	input.DisplayLetterNumber = normalizeOptionalString(input.DisplayLetterNumber)

	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		if desc == "" {
			input.Description = nil
		} else {
			input.Description = &desc
		}
	}

	return input
}

func validateCreateLetterInput(input CreateLetterInput) error {
	if input.ActorUserID == "" {
		return errors.New("actor user id is required")
	}

	if input.Title == "" {
		return errors.New("title is required")
	}

	if input.LetterDate == "" {
		return errors.New("letter_date is required")
	}

	if input.RegistrarName == "" {
		return errors.New("registrar_name is required")
	}

	if input.Sender == "" {
		return errors.New("sender is required")
	}

	if input.Receiver == "" {
		return errors.New("receiver is required")
	}

	return nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func resolveJalaliYear(letterDate time.Time, cfg LetterNumberConfig) int {
	if cfg.YearSource == "created_at" {
		now := time.Now().UTC()
		jy, _, _ := dateutil.GregorianToJalali(now.Year(), int(now.Month()), now.Day())
		return jy
	}

	jy, _, _ := dateutil.GregorianToJalali(letterDate.Year(), int(letterDate.Month()), letterDate.Day())
	return jy
}
