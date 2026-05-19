package imports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
	"github.com/hajimohammadinet/dabir/internal/domain/importjob"
)

type PreviewLettersImportUseCase struct {
	importRepo  importjob.Repository
	parser      *LetterExcelParser
	auditLogger *auditapp.Logger
}

type PreviewLettersImportInput struct {
	FileName    string
	FileReader  io.Reader
	ActorUserID string
	IPAddress   *string
	UserAgent   *string
}

func NewPreviewLettersImportUseCase(
	importRepo importjob.Repository,
	parser *LetterExcelParser,
	auditLogger *auditapp.Logger,
) *PreviewLettersImportUseCase {
	return &PreviewLettersImportUseCase{
		importRepo:  importRepo,
		parser:      parser,
		auditLogger: auditLogger,
	}
}

func (uc *PreviewLettersImportUseCase) Execute(ctx context.Context, input PreviewLettersImportInput) (*ImportJobDTO, error) {
	if input.ActorUserID == "" {
		return nil, fmt.Errorf("actor user id is required")
	}

	result, err := uc.parser.Parse(input.FileName, input.FileReader)
	if err != nil {
		return nil, err
	}

	detectedColumns, err := json.Marshal(result.DetectedColumns)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal detected columns: %w", err)
	}

	previewData, err := json.Marshal(result.Rows)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preview data: %w", err)
	}

	errorsData, err := json.Marshal(result.Errors)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal import errors: %w", err)
	}

	job := &importjob.ImportJob{
		Type:            importjob.TypeLetters,
		Status:          importjob.StatusPreviewed,
		FileName:        input.FileName,
		TotalRows:       result.TotalRows,
		ValidRows:       result.ValidRows,
		InvalidRows:     result.InvalidRows,
		MaxLetterNumber: result.MaxLetterNumber,
		DetectedColumns: detectedColumns,
		PreviewData:     previewData,
		Errors:          errorsData,
		CreatedBy:       input.ActorUserID,
	}

	if err := uc.importRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	if uc.auditLogger != nil {
		actorID := input.ActorUserID
		entityID := job.ID

		uc.auditLogger.Log(ctx, auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionLettersImportPreviewed,
			EntityType:  "import_job",
			EntityID:    &entityID,
			IPAddress:   input.IPAddress,
			UserAgent:   input.UserAgent,
			NewValue: map[string]interface{}{
				"file_name":         job.FileName,
				"total_rows":        job.TotalRows,
				"valid_rows":        job.ValidRows,
				"invalid_rows":      job.InvalidRows,
				"max_letter_number": job.MaxLetterNumber,
			},
		})
	}

	dto := ToImportJobDTO(*job)

	return &dto, nil
}
