package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	authapp "github.com/hajimohammadinet/dabir/internal/application/auth"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
	domainaudit "github.com/hajimohammadinet/dabir/internal/domain/audit"
)

type AuthHandler struct {
	loginUseCase          *authapp.LoginUseCase
	meUseCase             *authapp.MeUseCase
	changePasswordUseCase *authapp.ChangePasswordUseCase
	auditLogger           *auditapp.Logger
}

func NewAuthHandler(
	loginUseCase *authapp.LoginUseCase,
	meUseCase *authapp.MeUseCase,
	changePasswordUseCase *authapp.ChangePasswordUseCase,
	auditLogger *auditapp.Logger,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase:          loginUseCase,
		meUseCase:             meUseCase,
		changePasswordUseCase: changePasswordUseCase,
		auditLogger:           auditLogger,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input authapp.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.IPAddress = requestIP(r)
	input.UserAgent = requestUserAgent(r)

	output, err := h.loginUseCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, authapp.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
			return
		}

		if errors.Is(err, authapp.ErrInactiveUser) {
			response.Error(w, http.StatusForbidden, "INACTIVE_USER", "user is inactive")
			return
		}

		response.Error(w, http.StatusInternalServerError, "LOGIN_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	output, err := h.meUseCase.Execute(r.Context(), authUser.ID)
	if err != nil {
		if errors.Is(err, authapp.ErrUserNotFound) {
			response.Error(w, http.StatusUnauthorized, "USER_NOT_FOUND", "user not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "ME_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication is required")
		return
	}

	var input authapp.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	input.UserID = authUser.ID

	err := h.changePasswordUseCase.Execute(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, authapp.ErrInvalidCurrentPassword):
			response.Error(w, http.StatusBadRequest, "INVALID_CURRENT_PASSWORD", "invalid current password")
		case errors.Is(err, authapp.ErrWeakPassword):
			response.Error(w, http.StatusBadRequest, "WEAK_PASSWORD", "new password must be at least 8 characters")
		case errors.Is(err, authapp.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
		default:
			response.Error(w, http.StatusBadRequest, "CHANGE_PASSWORD_FAILED", err.Error())
		}
		return
	}

	if h.auditLogger != nil {
		actorID := authUser.ID

		h.auditLogger.Log(r.Context(), auditapp.LogInput{
			ActorUserID: &actorID,
			Action:      domainaudit.ActionUserPasswordChanged,
			EntityType:  "user",
			EntityID:    &actorID,
			IPAddress:   requestIP(r),
			UserAgent:   requestUserAgent(r),
			NewValue: map[string]interface{}{
				"changed": true,
			},
		})
	}

	response.JSON(w, http.StatusOK, map[string]bool{
		"changed": true,
	})
}
