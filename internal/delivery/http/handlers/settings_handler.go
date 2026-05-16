package handlers

import (
	"net/http"

	settingsapp "github.com/hajimohammadinet/dabir/internal/application/settings"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
)

type SettingsHandler struct {
	getPublicSettingsUseCase *settingsapp.GetPublicSettingsUseCase
}

func NewSettingsHandler(
	getPublicSettingsUseCase *settingsapp.GetPublicSettingsUseCase,
) *SettingsHandler {
	return &SettingsHandler{
		getPublicSettingsUseCase: getPublicSettingsUseCase,
	}
}

func (h *SettingsHandler) Public(w http.ResponseWriter, r *http.Request) {
	output, err := h.getPublicSettingsUseCase.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "GET_SETTINGS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}
