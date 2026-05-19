package imports

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
	"github.com/hajimohammadinet/dabir/internal/domain/importjob"
	"github.com/hajimohammadinet/dabir/internal/domain/letter"
)

var (
	ErrImportHasInvalidRows   = errors.New("import job has invalid rows")
	ErrImportAlreadyCommitted = errors.New("import job is already committed")
)

type CommitLettersImportUseCase struct {
	importRepo  importjob.Repository
	letterRepo  letter.Repository
	auditLogger *auditapp.Logger
}

type CommitLettersImportInput struct {
	ImportJobID   string
	ActorUserID   string
	RegistrarName string
	IPAddress     *string
	UserAgent     *string
}

type CommitLettersImportOutput struct {
	ImportID         string `json:"import_id"`
	ImportedRows     int    `json:"imported_rows"`
	SkippedRows      int    `json:"skipped_rows"`
	NextLetterNumber int64  `json:"next_letter_number"`
}

func NewCommitLettersImportUseCase(
	importRepo importjob.Repository,
	letterRepo letter.Repository,
	auditLogger *auditapp.Logger,
) *CommitLettersImportUseCase {
	return &CommitLettersImportUseCase{
		importRepo:  importRepo,
		letterRepo:  letterRepo,
		auditLogger: auditLogger,
	}
}

func (uc *CommitLettersImportUseCase) Execute(ctx context.Context, input CommitLettersImportInput) (*CommitLettersImportOutput, error) {
	job, err := uc.importRepo.FindByID(ctx, input.ImportJobID)
	if err != nil {
		return nil, fmt.Errorf("failed to find import job: %w", err)
	}

	if job == nil {
		return nil, ErrImportJobNotFound
	}

	if job.Status == importjob.StatusCommitted {
		return nil, ErrImportAlreadyCommitted
	}

	if job.InvalidRows > 0 {
		return nil, ErrImportHasInvalidRows
	}

	rows := make([]ImportedLetterRow, 0)
	if err := json.Unmarshal(job.PreviewData, &rows); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import preview data: %w", err)
	}

	letters := make([]letter.Letter, 0, len(rows))
	for _, row := range rows {
		letterDate, err := time.Parse("2006-01-02", row.LetterDate)
		if err != nil {
			return nil, fmt.Errorf("invalid parsed letter date at row %d: %w", row.RowNumber, err)
		}

		letters = append(letters, letter.Letter{
			LetterNumber:  row.LetterNumber,
			Title:         row.Title,
			LetterDate:    letterDate,
			RegistrarName: input.RegistrarName,
			Sender:        row.Sender,
			Receiver:      row.Receiver,
			CreatedBy:     input.ActorUserID,
			IsDeleted:     false,
		})
	}

	if err := uc.letterRepo.BulkCreate(ctx, letters); err != nil {
		return nil, err
	}

	if job.MaxLetterNumber != nil {
		if err := uc.letterRepo.SetSequenceValue(ctx, *job.MaxLetterNumber); err != nil {
			return nil, err
		}
	}

	if err := uc.importRepo.MarkCommitted(ctx, job.ID, input.ActorUserID); err != nil {
		return nil, err
	}

	nextLetterNumber := int64(1)
	if job.MaxLetterNumber != nil {
		nextLetterNumber = *job.MaxLetterNumber + 1
	}

	output := &CommitLettersImportOutput{
		ImportID:         job.ID,
		ImportedRows:     len(rows),
		SkippedRows:      job.InvalidRows,
		NextLetterNumber: nextLetterNumber,
	}

	if uc.auditLogger != nil {
		actorID := input.ActorUserID
		entityID := job.ID

		uc.auditLogger.Log(ctx, auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionLettersImportCommitted,
			EntityType:  "import_job",
			EntityID:    &entityID,
			IPAddress:   input.IPAddress,
			UserAgent:   input.UserAgent,
			NewValue:    output,
		})
	}

	return output, nil
}
