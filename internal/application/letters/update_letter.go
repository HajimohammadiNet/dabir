package letters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type UpdateLetterUseCase struct {
	letterRepo     letter.Repository
	configProvider *LetterConfigProvider
}

type UpdateLetterInput struct {
	ID string `json:"-"`

	Title      string `json:"title"`
	LetterDate string `json:"letter_date"`

	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`

	Description *string `json:"description"`

	ActorUserID string `json:"-"`
}

func NewUpdateLetterUseCase(
	letterRepo letter.Repository,
	configProvider *LetterConfigProvider,
) *UpdateLetterUseCase {
	return &UpdateLetterUseCase{
		letterRepo:     letterRepo,
		configProvider: configProvider,
	}
}

func (uc *UpdateLetterUseCase) Execute(ctx context.Context, input UpdateLetterInput) (*LetterDTO, error) {
	input = normalizeUpdateLetterInput(input)

	if err := validateUpdateLetterInput(input); err != nil {
		return nil, err
	}

	letterDate, err := time.Parse("2006-01-02", input.LetterDate)
	if err != nil {
		return nil, errors.New("letter_date must be in YYYY-MM-DD format")
	}

	l, err := uc.letterRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find letter: %w", err)
	}

	if l == nil || l.IsDeleted {
		return nil, ErrLetterNotFound
	}

	l.Title = input.Title
	l.LetterDate = letterDate
	l.Sender = input.Sender
	l.Receiver = input.Receiver
	l.Description = input.Description
	l.UpdatedBy = &input.ActorUserID

	if err := uc.letterRepo.Update(ctx, l); err != nil {
		return nil, fmt.Errorf("failed to update letter: %w", err)
	}

	cfg := uc.configProvider.Get(ctx)
	dto := ToLetterDTO(*l, cfg)

	return &dto, nil
}

func normalizeUpdateLetterInput(input UpdateLetterInput) UpdateLetterInput {
	input.Title = strings.TrimSpace(input.Title)
	input.LetterDate = strings.TrimSpace(input.LetterDate)
	input.Sender = strings.TrimSpace(input.Sender)
	input.Receiver = strings.TrimSpace(input.Receiver)

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

func validateUpdateLetterInput(input UpdateLetterInput) error {
	if input.ID == "" {
		return errors.New("letter id is required")
	}

	if input.ActorUserID == "" {
		return errors.New("actor user id is required")
	}

	if input.Title == "" {
		return errors.New("title is required")
	}

	if input.LetterDate == "" {
		return errors.New("letter_date is required")
	}

	if input.Sender == "" {
		return errors.New("sender is required")
	}

	if input.Receiver == "" {
		return errors.New("receiver is required")
	}

	return nil
}
