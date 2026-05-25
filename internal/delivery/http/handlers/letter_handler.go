package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	lettersapp "github.com/hajimohammadinet/dabir/internal/application/letters"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type LetterHandler struct {
	createLetterUseCase *lettersapp.CreateLetterUseCase
	listLettersUseCase  *lettersapp.ListLettersUseCase
	getLetterUseCase    *lettersapp.GetLetterUseCase
	updateLetterUseCase *lettersapp.UpdateLetterUseCase
	deleteLetterUseCase *lettersapp.DeleteLetterUseCase
	auditLogger         *auditapp.Logger
}

func NewLetterHandler(
	createLetterUseCase *lettersapp.CreateLetterUseCase,
	listLettersUseCase *lettersapp.ListLettersUseCase,
	getLetterUseCase *lettersapp.GetLetterUseCase,
	updateLetterUseCase *lettersapp.UpdateLetterUseCase,
	deleteLetterUseCase *lettersapp.DeleteLetterUseCase,
	auditLogger *auditapp.Logger,
) *LetterHandler {
	return &LetterHandler{
		createLetterUseCase: createLetterUseCase,
		listLettersUseCase:  listLettersUseCase,
		getLetterUseCase:    getLetterUseCase,
		updateLetterUseCase: updateLetterUseCase,
		deleteLetterUseCase: deleteLetterUseCase,
		auditLogger:         auditLogger,
	}
}

func (h *LetterHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	var input lettersapp.CreateLetterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.ActorUserID = authUser.ID
	input.RegistrarName = authUser.Username

	output, err := h.createLetterUseCase.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_LETTER_FAILED", err.Error())
		return
	}

	actorID := authUser.ID
	entityID := output.ID

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionLetterCreated,
		EntityType:  "letter",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		NewValue:    output,
	})

	response.JSON(w, http.StatusCreated, output)
}

func (h *LetterHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := parseIntQuery(query.Get("page"), 1)
	pageSize := parseIntQuery(query.Get("page_size"), 20)

	includeDeleted := false
	if value := strings.TrimSpace(query.Get("include_deleted")); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "INVALID_INCLUDE_DELETED", "include_deleted must be true or false")
			return
		}
		includeDeleted = parsed
	}

	output, err := h.listLettersUseCase.Execute(r.Context(), lettersapp.ListLettersInput{
		Page:           page,
		PageSize:       pageSize,
		Search:         query.Get("search"),
		RegistrarName:  query.Get("registrar_name"),
		Sender:         query.Get("sender"),
		Receiver:       query.Get("receiver"),
		FromDate:       query.Get("from_date"),
		ToDate:         query.Get("to_date"),
		IncludeDeleted: includeDeleted,
		SortBy:         query.Get("sort_by"),
		SortOrder:      query.Get("sort_order"),
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, "LIST_LETTERS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *LetterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.getLetterUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, lettersapp.ErrLetterNotFound) {
			response.Error(w, http.StatusNotFound, "LETTER_NOT_FOUND", "letter not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "GET_LETTER_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *LetterHandler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	id := chi.URLParam(r, "id")

	oldOutput, _ := h.getLetterUseCase.Execute(r.Context(), id)

	var input lettersapp.UpdateLetterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.ID = id
	input.ActorUserID = authUser.ID

	output, err := h.updateLetterUseCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, lettersapp.ErrLetterNotFound) {
			response.Error(w, http.StatusNotFound, "LETTER_NOT_FOUND", "letter not found")
			return
		}

		response.Error(w, http.StatusBadRequest, "UPDATE_LETTER_FAILED", err.Error())
		return
	}

	actorID := authUser.ID
	entityID := output.ID

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionLetterUpdated,
		EntityType:  "letter",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		OldValue:    oldOutput,
		NewValue:    output,
	})

	response.JSON(w, http.StatusOK, output)
}

func (h *LetterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())

	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	id := chi.URLParam(r, "id")
	oldOutput, _ := h.getLetterUseCase.Execute(r.Context(), id)

	err := h.deleteLetterUseCase.Execute(r.Context(), lettersapp.DeleteLetterInput{
		ID:          id,
		ActorUserID: authUser.ID,
	})
	if err != nil {
		if errors.Is(err, lettersapp.ErrLetterNotFound) {
			response.Error(w, http.StatusNotFound, "LETTER_NOT_FOUND", "letter not found")
			return
		}

		response.Error(w, http.StatusBadRequest, "DELETE_LETTER_FAILED", err.Error())
		return
	}

	actorID := authUser.ID
	entityID := id

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionLetterDeleted,
		EntityType:  "letter",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		OldValue:    oldOutput,
		NewValue: map[string]interface{}{
			"deleted": true,
		},
	})

	response.JSON(w, http.StatusOK, map[string]bool{
		"deleted": true,
	})
}
