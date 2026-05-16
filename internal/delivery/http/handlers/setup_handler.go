package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	setupapp "github.com/hajimohammadinet/dabir/internal/application/setup"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
)

type SetupHandler struct {
	checkStatusUseCase *setupapp.CheckStatusUseCase
	initializeUseCase  *setupapp.InitializeUseCase
}

func NewSetupHandler(
	checkStatusUseCase *setupapp.CheckStatusUseCase,
	initializeUseCase *setupapp.InitializeUseCase,
) *SetupHandler {
	return &SetupHandler{
		checkStatusUseCase: checkStatusUseCase,
		initializeUseCase:  initializeUseCase,
	}
}

func (h *SetupHandler) Status(w http.ResponseWriter, r *http.Request) {
	output, err := h.checkStatusUseCase.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "SETUP_STATUS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}

func (h *SetupHandler) Initialize(w http.ResponseWriter, r *http.Request) {
	var input setupapp.InitializeInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	output, err := h.initializeUseCase.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, setupapp.ErrAlreadyInitialized) {
			response.Error(w, http.StatusConflict, "ALREADY_INITIALIZED", "application is already initialized")
			return
		}

		response.Error(w, http.StatusBadRequest, "INITIALIZE_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, output)
}
