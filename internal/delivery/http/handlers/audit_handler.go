package handlers

import (
	"net/http"

	auditapp "github.com/hajimohammadinet/dabir/internal/application/audit"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
)

type AuditHandler struct {
	listAuditLogsUseCase *auditapp.ListAuditLogsUseCase
}

func NewAuditHandler(listAuditLogsUseCase *auditapp.ListAuditLogsUseCase) *AuditHandler {
	return &AuditHandler{
		listAuditLogsUseCase: listAuditLogsUseCase,
	}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	output, err := h.listAuditLogsUseCase.Execute(r.Context(), auditapp.ListAuditLogsInput{
		Page:        parseIntQuery(query.Get("page"), 1),
		PageSize:    parseIntQuery(query.Get("page_size"), 20),
		Action:      query.Get("action"),
		EntityType:  query.Get("entity_type"),
		ActorUserID: query.Get("actor_user_id"),
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "LIST_AUDIT_LOGS_FAILED", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, output)
}
