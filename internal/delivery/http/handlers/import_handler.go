package handlers

import (
	"errors"
	"net/http"

	importsapp "github.com/hajimohammadinet/dabir/internal/application/imports"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"

	"github.com/go-chi/chi/v5"
)

type ImportHandler struct {
	previewLettersImportUseCase *importsapp.PreviewLettersImportUseCase
	commitLettersImportUseCase  *importsapp.CommitLettersImportUseCase
	getImportJobUseCase         *importsapp.GetImportJobUseCase
}

func NewImportHandler(
	previewLettersImportUseCase *importsapp.PreviewLettersImportUseCase,
	commitLettersImportUseCase *importsapp.CommitLettersImportUseCase,
	getImportJobUseCase *importsapp.GetImportJobUseCase,
) *ImportHandler {
	return &ImportHandler{
		previewLettersImportUseCase: previewLettersImportUseCase,
		commitLettersImportUseCase:  commitLettersImportUseCase,
		getImportJobUseCase:         getImportJobUseCase,
	}
}

func (h *ImportHandler) PreviewLetters(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_MULTIPART_FORM", "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_REQUIRED", "file field is required")
		return
	}
	defer file.Close()

	output, err := h.previewLettersImportUseCase.Execute(r.Context(), importsapp.PreviewLettersImportInput{
		FileName:    header.Filename,
		FileReader:  file,
		ActorUserID: authUser.ID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "IMPORT_PREVIEW_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, output)
}

func (h *ImportHandler) CommitLetters(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	importID := chi.URLParam(r, "id")

	output, err := h.commitLettersImportUseCase.Execute(r.Context(), importsapp.CommitLettersImportInput{
		ImportJobID:   importID,
		ActorUserID:   authUser.ID,
		RegistrarName: authUser.Username,
		IPAddress:     requestIP(r),
		UserAgent:     requestUserAgent(r),
	})
	if err != nil {
		switch {
		case errors.Is(err, importsapp.ErrImportJobNotFound):
			response.Error(w, http.StatusNotFound, "IMPORT_JOB_NOT_FOUND", "import job not found")
		case errors.Is(err, importsapp.ErrImportHasInvalidRows):
			response.Error(w, http.StatusBadRequest, "IMPORT_HAS_INVALID_ROWS", "import job has invalid rows")
		case errors.Is(err, importsapp.ErrImportAlreadyCommitted):
			response.Error(w, http.StatusConflict, "IMPORT_ALREADY_COMMITTED", "import job is already committed")
		default:
			response.Error(w, http.StatusBadRequest, "IMPORT_COMMIT_FAILED", err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *ImportHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	importID := chi.URLParam(r, "id")

	output, err := h.getImportJobUseCase.Execute(r.Context(), importID)
	if err != nil {
		if errors.Is(err, importsapp.ErrImportJobNotFound) {
			response.Error(w, http.StatusNotFound, "IMPORT_JOB_NOT_FOUND", "import job not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "GET_IMPORT_JOB_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}
