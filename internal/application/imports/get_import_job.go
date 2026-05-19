package imports

import (
	"context"
	"errors"
	"fmt"

	"github.com/hajimohammadinet/dabir/internal/domain/importjob"
)

var ErrImportJobNotFound = errors.New("import job not found")

type GetImportJobUseCase struct {
	importRepo importjob.Repository
}

func NewGetImportJobUseCase(importRepo importjob.Repository) *GetImportJobUseCase {
	return &GetImportJobUseCase{
		importRepo: importRepo,
	}
}

func (uc *GetImportJobUseCase) Execute(ctx context.Context, id string) (*ImportJobDTO, error) {
	job, err := uc.importRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find import job: %w", err)
	}

	if job == nil {
		return nil, ErrImportJobNotFound
	}

	dto := ToImportJobDTO(*job)

	return &dto, nil
}
