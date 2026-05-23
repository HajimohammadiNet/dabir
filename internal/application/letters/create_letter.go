package letters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	//"time"

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

	nextNumber, err := uc.letterRepo.NextNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate letter number: %w", err)
	}

	l := &letter.Letter{
		LetterNumber: nextNumber,
		Title:        input.Title,
		LetterDate:   letterDate,

		RegistrarName: input.RegistrarName,
		Sender:        input.Sender,
		Receiver:      input.Receiver,

		Description: normalizeOptionalString(input.Description),

		CreatedBy: input.ActorUserID,
		IsDeleted: false,
	}

	if err := uc.letterRepo.Create(ctx, l); err != nil {
		return nil, fmt.Errorf("failed to create letter: %w", err)
	}

	cfg := uc.configProvider.Get(ctx)
	dto := ToLetterDTO(*l, cfg)

	return &dto, nil
}

func normalizeCreateLetterInput(input CreateLetterInput) CreateLetterInput {
	input.Title = strings.TrimSpace(input.Title)
	input.LetterDate = strings.TrimSpace(input.LetterDate)
	input.Sender = strings.TrimSpace(input.Sender)
	input.Receiver = strings.TrimSpace(input.Receiver)
	input.RegistrarName = strings.TrimSpace(input.RegistrarName)

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
