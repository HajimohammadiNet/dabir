package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	settingsapp "github.com/hajimohammadinet/dabir/internal/application/settings"
	setupapp "github.com/hajimohammadinet/dabir/internal/application/setup"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/handlers"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/postgres"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/security"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)

	userRepo := postgres.NewUserRepository(db)
	settingsRepo := postgres.NewSettingsRepository(db)
	passwordHasher := security.NewPasswordHasher()

	checkStatusUseCase := setupapp.NewCheckStatusUseCase(userRepo)
	initializeUseCase := setupapp.NewInitializeUseCase(
		userRepo,
		settingsRepo,
		passwordHasher,
	)

	setupHandler := handlers.NewSetupHandler(
		checkStatusUseCase,
		initializeUseCase,
	)

	getPublicSettingsUseCase := settingsapp.NewGetPublicSettingsUseCase(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(getPublicSettingsUseCase)

	r.Get("/healthz", healthHandler.Healthz)
	r.Get("/readyz", healthHandler.Readyz)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", healthHandler.Healthz)
		r.Get("/readyz", healthHandler.Readyz)

		r.Route("/setup", func(r chi.Router) {
			r.Get("/status", setupHandler.Status)
			r.Post("/initialize", setupHandler.Initialize)
		})

		r.Route("/settings", func(r chi.Router) {
			r.Get("/public", settingsHandler.Public)
		})
	})

	return r
}
