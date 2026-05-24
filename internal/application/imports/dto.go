package imports

import (
	"encoding/json"
	"time"

	"github.com/hajimohammadinet/dabir/internal/domain/importjob"
)

type ImportErrorDTO struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ImportedLetterRow struct {
	RowNumber int `json:"row_number"`

	LetterNumber        int64  `json:"letter_number"`
	DisplayLetterNumber string `json:"display_letter_number"`

	Title            string `json:"title"`
	LetterDate       string `json:"letter_date"`
	LetterDateJalali string `json:"letter_date_jalali"`
	Sender           string `json:"sender"`
	Receiver         string `json:"receiver"`
}

type ImportJobDTO struct {
	ID string `json:"id"`

	Type   importjob.Type   `json:"type"`
	Status importjob.Status `json:"status"`

	FileName string `json:"file_name"`

	TotalRows   int `json:"total_rows"`
	ValidRows   int `json:"valid_rows"`
	InvalidRows int `json:"invalid_rows"`

	MaxLetterNumber *int64 `json:"max_letter_number,omitempty"`

	DetectedColumns map[string]string   `json:"detected_columns,omitempty"`
	PreviewData     []ImportedLetterRow `json:"preview_data,omitempty"`
	Errors          []ImportErrorDTO    `json:"errors,omitempty"`

	CreatedBy   string  `json:"created_by"`
	CommittedBy *string `json:"committed_by,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	CommittedAt *time.Time `json:"committed_at,omitempty"`
}

func ToImportJobDTO(job importjob.ImportJob) ImportJobDTO {
	return ImportJobDTO{
		ID:              job.ID,
		Type:            job.Type,
		Status:          job.Status,
		FileName:        job.FileName,
		TotalRows:       job.TotalRows,
		ValidRows:       job.ValidRows,
		InvalidRows:     job.InvalidRows,
		MaxLetterNumber: job.MaxLetterNumber,
		DetectedColumns: decodeMap(job.DetectedColumns),
		PreviewData:     decodeRows(job.PreviewData),
		Errors:          decodeErrors(job.Errors),
		CreatedBy:       job.CreatedBy,
		CommittedBy:     job.CommittedBy,
		CreatedAt:       job.CreatedAt,
		CommittedAt:     job.CommittedAt,
	}
}

func decodeMap(raw []byte) map[string]string {
	if len(raw) == 0 {
		return nil
	}

	var value map[string]string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}

	return value
}

func decodeRows(raw []byte) []ImportedLetterRow {
	if len(raw) == 0 {
		return nil
	}

	var value []ImportedLetterRow
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}

	return value
}

func decodeErrors(raw []byte) []ImportErrorDTO {
	if len(raw) == 0 {
		return nil
	}

	var value []ImportErrorDTO
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}

	return value
}
