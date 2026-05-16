package letters

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type ListLettersUseCase struct {
	letterRepo     letter.Repository
	configProvider *LetterConfigProvider
}

type ListLettersInput struct {
	Page     int
	PageSize int

	Search        string
	Destination   string
	RegistrarName string

	FromDate string
	ToDate   string

	IncludeDeleted bool
}

type ListLettersOutput struct {
	Items      []LetterDTO `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func NewListLettersUseCase(
	letterRepo letter.Repository,
	configProvider *LetterConfigProvider,
) *ListLettersUseCase {
	return &ListLettersUseCase{
		letterRepo:     letterRepo,
		configProvider: configProvider,
	}
}

func (uc *ListLettersUseCase) Execute(ctx context.Context, input ListLettersInput) (*ListLettersOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	input.Search = strings.TrimSpace(input.Search)
	input.Destination = strings.TrimSpace(input.Destination)
	input.RegistrarName = strings.TrimSpace(input.RegistrarName)

	filter := letter.ListFilter{
		Page:           input.Page,
		PageSize:       input.PageSize,
		Search:         input.Search,
		Destination:    input.Destination,
		RegistrarName:  input.RegistrarName,
		IncludeDeleted: input.IncludeDeleted,
	}

	if input.FromDate != "" {
		fromDate, err := parseDate(input.FromDate, "from_date")
		if err != nil {
			return nil, err
		}

		filter.FromDate = &fromDate
	}

	if input.ToDate != "" {
		toDate, err := parseDate(input.ToDate, "to_date")
		if err != nil {
			return nil, err
		}

		filter.ToDate = &toDate
	}

	lettersList, total, err := uc.letterRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	cfg := uc.configProvider.Get(ctx)

	items := make([]LetterDTO, 0, len(lettersList))
	for _, l := range lettersList {
		items = append(items, ToLetterDTO(l, cfg))
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + input.PageSize - 1) / input.PageSize
	}

	return &ListLettersOutput{
		Items:      items,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}

func parseDate(value string, fieldName string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, errors.New(fieldName + " must be in YYYY-MM-DD format")
	}

	return parsed, nil
}
