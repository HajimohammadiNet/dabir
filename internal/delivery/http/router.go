package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hajimohammadinet/dabir/internal/delivery/http/handlers"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)

	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", healthHandler.Healthz)
		r.Get("/readyz", healthHandler.Readyz)
	})

	return r
}
