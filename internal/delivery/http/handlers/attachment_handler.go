package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	attachmentsapp "github.com/hajimohammadinet/dabir/internal/application/attachments"
	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	lettersapp "github.com/hajimohammadinet/dabir/internal/application/letters"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type AttachmentHandler struct {
	uploadAttachmentUseCase         *attachmentsapp.UploadAttachmentUseCase
	listAttachmentsUseCase          *attachmentsapp.ListAttachmentsUseCase
	getAttachmentDownloadURLUseCase *attachmentsapp.GetAttachmentDownloadURLUseCase
	deleteAttachmentUseCase         *attachmentsapp.DeleteAttachmentUseCase
	getLetterUseCase                *lettersapp.GetLetterUseCase
	auditLogger                     *auditapp.Logger
}

func NewAttachmentHandler(
	uploadAttachmentUseCase *attachmentsapp.UploadAttachmentUseCase,
	listAttachmentsUseCase *attachmentsapp.ListAttachmentsUseCase,
	getAttachmentDownloadURLUseCase *attachmentsapp.GetAttachmentDownloadURLUseCase,
	deleteAttachmentUseCase *attachmentsapp.DeleteAttachmentUseCase,
	getLetterUseCase *lettersapp.GetLetterUseCase,
	auditLogger *auditapp.Logger,
) *AttachmentHandler {
	return &AttachmentHandler{
		uploadAttachmentUseCase:         uploadAttachmentUseCase,
		listAttachmentsUseCase:          listAttachmentsUseCase,
		getAttachmentDownloadURLUseCase: getAttachmentDownloadURLUseCase,
		deleteAttachmentUseCase:         deleteAttachmentUseCase,
		getLetterUseCase:                getLetterUseCase,
		auditLogger:                     auditLogger,
	}
}

func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	letterID := chi.URLParam(r, "id")
	letterItem, _ := h.getLetterUseCase.Execute(r.Context(), letterID)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_MULTIPART_FORM", "invalid multipart form")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		response.Error(w, http.StatusBadRequest, "NO_FILES", "no files uploaded")
		return
	}

	result := make([]attachmentsapp.AttachmentDTO, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			response.Error(w, http.StatusBadRequest, "INVALID_FILE", err.Error())
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType(make([]byte, 512))
		}

		dto, err := h.uploadAttachmentUseCase.Execute(r.Context(), attachmentsapp.UploadAttachmentInput{
			LetterID:    letterID,
			FileName:    fileHeader.Filename,
			ContentType: contentType,
			SizeBytes:   fileHeader.Size,
			Reader:      file,
			ActorUserID: authUser.ID,
		})
		_ = file.Close()

		if err != nil {
			response.Error(w, http.StatusBadRequest, "UPLOAD_ATTACHMENT_FAILED", err.Error())
			return
		}

		result = append(result, *dto)

		if letterItem != nil {
			actorID := authUser.ID
			entityID := letterID
			h.auditLogger.Log(r.Context(), auditapp.LogInput{
				ActorUserID: &actorID,
				Action:      domainaudit.ActionLetterAttachmentUploaded,
				EntityType:  "letter",
				EntityID:    &entityID,
				IPAddress:   requestIP(r),
				UserAgent:   requestUserAgent(r),
				NewValue: map[string]interface{}{
					"direction":  letterItem.Direction,
					"attachment": dto,
				},
			})
		}
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	letterID := chi.URLParam(r, "id")

	output, err := h.listAttachmentsUseCase.Execute(r.Context(), attachmentsapp.ListAttachmentsInput{
		LetterID: letterID,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "LIST_ATTACHMENTS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *AttachmentHandler) DownloadURL(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	letterID := chi.URLParam(r, "id")
	attachmentID := chi.URLParam(r, "attachment_id")

	output, err := h.getAttachmentDownloadURLUseCase.Execute(r.Context(), attachmentsapp.GetAttachmentDownloadURLInput{
		LetterID:     letterID,
		AttachmentID: attachmentID,
		ActorUserID:  authUser.ID,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "GET_ATTACHMENT_DOWNLOAD_URL_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *AttachmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	letterID := chi.URLParam(r, "id")
	attachmentID := chi.URLParam(r, "attachment_id")
	letterItem, _ := h.getLetterUseCase.Execute(r.Context(), letterID)

	deletedAttachment, err := h.deleteAttachmentUseCase.Execute(r.Context(), attachmentsapp.DeleteAttachmentInput{
		LetterID:     letterID,
		AttachmentID: attachmentID,
		ActorUserID:  authUser.ID,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "DELETE_ATTACHMENT_FAILED", err.Error())
		return
	}

	if letterItem != nil {
		actorID := authUser.ID
		entityID := letterID
		h.auditLogger.Log(r.Context(), auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionLetterAttachmentDeleted,
			EntityType:  "letter",
			EntityID:    &entityID,
			IPAddress:   requestIP(r),
			UserAgent:   requestUserAgent(r),
			OldValue: map[string]interface{}{
				"direction":  letterItem.Direction,
				"attachment": deletedAttachment,
			},
			NewValue: map[string]interface{}{
				"direction":     letterItem.Direction,
				"attachment_id": attachmentID,
				"deleted":       true,
			},
		})
	}

	response.JSON(w, http.StatusOK, map[string]bool{
		"deleted": true,
	})
}
