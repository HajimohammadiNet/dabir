package letters

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

type SuggestLetterNumberUseCase struct {
	letterRepo letter.Repository
}

type SuggestLetterNumberInput struct {
	Prefix    string
	Direction letter.Direction
}

type SuggestLetterNumberOutput struct {
	Prefix          string  `json:"prefix"`
	LastNumber      *string `json:"last_number"`
	SuggestedNumber *string `json:"suggested_number"`
}

func NewSuggestLetterNumberUseCase(letterRepo letter.Repository) *SuggestLetterNumberUseCase {
	return &SuggestLetterNumberUseCase{
		letterRepo: letterRepo,
	}
}

func (uc *SuggestLetterNumberUseCase) Execute(ctx context.Context, input SuggestLetterNumberInput) (*SuggestLetterNumberOutput, error) {
	if input.Direction == "" {
		input.Direction = letter.DirectionIncoming
	}
	if !input.Direction.IsValid() {
		return nil, fmt.Errorf("direction must be incoming or outgoing")
	}

	prefix := normalizePersianArabicDigits(strings.TrimSpace(input.Prefix))

	lastNumber, err := uc.letterRepo.FindLatestDisplayLetterNumberByPrefix(ctx, input.Direction, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to find latest letter number: %w", err)
	}

	var suggested *string

	if lastNumber != nil {
		value := suggestNextManualNumber(*lastNumber)
		if value != "" {
			suggested = &value
		}
	} else if prefix != "" {
		value := suggestInitialManualNumber(prefix)
		if value != "" {
			suggested = &value
		}
	}

	return &SuggestLetterNumberOutput{
		Prefix:          prefix,
		LastNumber:      lastNumber,
		SuggestedNumber: suggested,
	}, nil
}

func suggestNextManualNumber(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	normalized := normalizePersianArabicDigits(value)

	start, end := findLastDigitRange(normalized)
	if start == -1 || end == -1 {
		return ""
	}

	numericPart := normalized[start:end]
	number, err := strconv.Atoi(numericPart)
	if err != nil {
		return ""
	}

	nextNumber := strconv.Itoa(number + 1)
	if len(nextNumber) < len(numericPart) {
		nextNumber = strings.Repeat("0", len(numericPart)-len(nextNumber)) + nextNumber
	}

	return normalized[:start] + nextNumber + normalized[end:]
}

func suggestInitialManualNumber(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}

	normalized := normalizePersianArabicDigits(prefix)

	last := []rune(normalized)
	if len(last) == 0 {
		return ""
	}

	// If user typed a prefix ending with a separator, suggest 001.
	// Examples:
	// 405-ق- -> 405-ق-001
	// HR-2026- -> HR-2026-001
	if strings.HasSuffix(normalized, "-") ||
		strings.HasSuffix(normalized, "/") ||
		strings.HasSuffix(normalized, "_") ||
		strings.HasSuffix(normalized, ".") {
		return normalized + "001"
	}

	// If prefix already ends with a number, increment it.
	if unicode.IsDigit(last[len(last)-1]) {
		return suggestNextManualNumber(normalized)
	}

	return normalized
}

func findLastDigitRange(value string) (int, int) {
	end := -1

	for i := len(value) - 1; i >= 0; i-- {
		if value[i] >= '0' && value[i] <= '9' {
			end = i + 1
			break
		}
	}

	if end == -1 {
		return -1, -1
	}

	start := end - 1
	for start >= 0 && value[start] >= '0' && value[start] <= '9' {
		start--
	}

	return start + 1, end
}

func normalizePersianArabicDigits(value string) string {
	replacer := strings.NewReplacer(
		"۰", "0",
		"۱", "1",
		"۲", "2",
		"۳", "3",
		"۴", "4",
		"۵", "5",
		"۶", "6",
		"۷", "7",
		"۸", "8",
		"۹", "9",
		"٠", "0",
		"١", "1",
		"٢", "2",
		"٣", "3",
		"٤", "4",
		"٥", "5",
		"٦", "6",
		"٧", "7",
		"٨", "8",
		"٩", "9",
	)

	return replacer.Replace(value)
}
