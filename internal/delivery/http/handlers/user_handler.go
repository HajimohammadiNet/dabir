package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	usersapp "github.com/hajimohammadinet/dabir/internal/application/users"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
	"github.com/hajimohammadinet/dabir/internal/domain/user"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	createUserUseCase    *usersapp.CreateUserUseCase
	listUsersUseCase     *usersapp.ListUsersUseCase
	getUserUseCase       *usersapp.GetUserUseCase
	updateUserUseCase    *usersapp.UpdateUserUseCase
	setUserActiveUseCase *usersapp.SetUserActiveUseCase
	auditLogger          *auditapp.Logger
	resetPasswordUseCase *usersapp.ResetPasswordUseCase
}

func NewUserHandler(
	createUserUseCase *usersapp.CreateUserUseCase,
	listUsersUseCase *usersapp.ListUsersUseCase,
	getUserUseCase *usersapp.GetUserUseCase,
	updateUserUseCase *usersapp.UpdateUserUseCase,
	setUserActiveUseCase *usersapp.SetUserActiveUseCase,
	auditLogger *auditapp.Logger,
	resetPasswordUseCase *usersapp.ResetPasswordUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:    createUserUseCase,
		listUsersUseCase:     listUsersUseCase,
		getUserUseCase:       getUserUseCase,
		updateUserUseCase:    updateUserUseCase,
		setUserActiveUseCase: setUserActiveUseCase,
		auditLogger:          auditLogger,
		resetPasswordUseCase: resetPasswordUseCase,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input usersapp.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	output, err := h.createUserUseCase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usersapp.ErrUsernameAlreadyExists):
			response.Error(w, http.StatusConflict, "USERNAME_ALREADY_EXISTS", "username already exists")
		case errors.Is(err, usersapp.ErrInvalidRole):
			response.Error(w, http.StatusBadRequest, "INVALID_ROLE", "invalid role")
		default:
			response.Error(w, http.StatusBadRequest, "CREATE_USER_FAILED", err.Error())
		}
		return
	}

	authUser, _ := middleware.GetAuthUser(r.Context())
	actorID := authUser.ID
	entityID := output.ID

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionUserCreated,
		EntityType:  "user",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		NewValue:    output,
	})

	response.JSON(w, http.StatusCreated, output)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := parseIntQuery(query.Get("page"), 1)
	pageSize := parseIntQuery(query.Get("page_size"), 20)
	search := strings.TrimSpace(query.Get("search"))

	var roleFilter *user.Role
	if roleValue := strings.TrimSpace(query.Get("role")); roleValue != "" {
		role := user.Role(roleValue)
		roleFilter = &role
	}

	var isActiveFilter *bool
	if activeValue := strings.TrimSpace(query.Get("is_active")); activeValue != "" {
		active, err := strconv.ParseBool(activeValue)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "INVALID_IS_ACTIVE", "is_active must be true or false")
			return
		}

		isActiveFilter = &active
	}

	output, err := h.listUsersUseCase.Execute(r.Context(), usersapp.ListUsersInput{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Role:     roleFilter,
		IsActive: isActiveFilter,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "LIST_USERS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.getUserUseCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, usersapp.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "GET_USER_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	oldOutput, _ := h.getUserUseCase.Execute(r.Context(), id)

	var input usersapp.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.ID = id

	output, err := h.updateUserUseCase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usersapp.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		case errors.Is(err, usersapp.ErrInvalidRole):
			response.Error(w, http.StatusBadRequest, "INVALID_ROLE", "invalid role")
		default:
			response.Error(w, http.StatusBadRequest, "UPDATE_USER_FAILED", err.Error())
		}
		return
	}

	authUser, _ := middleware.GetAuthUser(r.Context())
	actorID := authUser.ID
	entityID := output.ID

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionUserUpdated,
		EntityType:  "user",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		OldValue:    oldOutput,
		NewValue:    output,
	})

	response.JSON(w, http.StatusOK, output)
}

func (h *UserHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.setUserActiveUseCase.Execute(r.Context(), id, false); err != nil {
		if errors.Is(err, usersapp.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}

		response.Error(w, http.StatusBadRequest, "DEACTIVATE_USER_FAILED", err.Error())
		return
	}

	authUser, _ := middleware.GetAuthUser(r.Context())
	actorID := authUser.ID
	entityID := id

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionUserDeactivated,
		EntityType:  "user",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		NewValue: map[string]interface{}{
			"is_active": false,
		},
	})

	response.JSON(w, http.StatusOK, map[string]bool{
		"deactivated": true,
	})
}

func (h *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.setUserActiveUseCase.Execute(r.Context(), id, true); err != nil {
		if errors.Is(err, usersapp.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
			return
		}

		response.Error(w, http.StatusBadRequest, "ACTIVATE_USER_FAILED", err.Error())
		return
	}

	authUser, _ := middleware.GetAuthUser(r.Context())
	actorID := authUser.ID
	entityID := id

	h.auditLogger.Log(r.Context(), auditapp.LogInput{
		ActorUserID: &actorID,
		Action:      domainaudit.ActionUserActivated,
		EntityType:  "user",
		EntityID:    &entityID,
		IPAddress:   requestIP(r),
		UserAgent:   requestUserAgent(r),
		NewValue: map[string]interface{}{
			"is_active": true,
		},
	})

	response.JSON(w, http.StatusOK, map[string]bool{
		"activated": true,
	})
}

func parseIntQuery(value string, defaultValue int) int {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	authUser, _ := middleware.GetAuthUser(r.Context())
	id := chi.URLParam(r, "id")

	var input usersapp.ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.UserID = id

	err := h.resetPasswordUseCase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usersapp.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		case errors.Is(err, usersapp.ErrWeakPassword):
			response.Error(w, http.StatusBadRequest, "WEAK_PASSWORD", "password must be at least 8 characters")
		default:
			response.Error(w, http.StatusBadRequest, "RESET_PASSWORD_FAILED", err.Error())
		}
		return
	}

	if h.auditLogger != nil {
		actorID := authUser.ID
		entityID := id

		h.auditLogger.Log(r.Context(), auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionUserPasswordReset,
			EntityType:  "user",
			EntityID:    &entityID,
			IPAddress:   requestIP(r),
			UserAgent:   requestUserAgent(r),
			NewValue: map[string]interface{}{
				"reset": true,
			},
		})
	}

	response.JSON(w, http.StatusOK, map[string]bool{
		"reset": true,
	})
}
