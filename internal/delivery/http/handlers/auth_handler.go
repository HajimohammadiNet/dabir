package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	authapp "github.com/hajimohammadinet/dabir/internal/application/auth"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
)

type AuthHandler struct {
	loginUseCase *authapp.LoginUseCase
	meUseCase    *authapp.MeUseCase
}

func NewAuthHandler(
	loginUseCase *authapp.LoginUseCase,
	meUseCase *authapp.MeUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase: loginUseCase,
		meUseCase:    meUseCase,
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
