package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	authapp "github.com/hajimohammadinet/dabir/internal/application/auth"
	settingsapp "github.com/hajimohammadinet/dabir/internal/application/settings"
	setupapp "github.com/hajimohammadinet/dabir/internal/application/setup"
	"github.com/hajimohammadinet/dabir/internal/config"
	"github.com/hajimohammadinet/dabir/internal/delivery/http/handlers"
	httpmiddleware "github.com/hajimohammadinet/dabir/internal/delivery/http/middleware"
	"github.com/hajimohammadinet/dabir/internal/domain/user"
	infraauth "github.com/hajimohammadinet/dabir/internal/infrastructure/auth"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/postgres"
	"github.com/hajimohammadinet/dabir/internal/infrastructure/security"
)

func NewRouter(db *pgxpool.Pool, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)

	userRepo := postgres.NewUserRepository(db)
	settingsRepo := postgres.NewSettingsRepository(db)
	passwordHasher := security.NewPasswordHasher()
	jwtService := infraauth.NewJWTService(cfg.Auth)

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

	loginUseCase := authapp.NewLoginUseCase(
		userRepo,
		passwordHasher,
		jwtService,
	)
	meUseCase := authapp.NewMeUseCase(userRepo)
	authHandler := handlers.NewAuthHandler(loginUseCase, meUseCase)

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

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)

			r.Group(func(r chi.Router) {
				r.Use(httpmiddleware.AuthMiddleware(jwtService))

				r.Get("/me", authHandler.Me)

				r.With(httpmiddleware.RequireRoles(user.RoleSuperUser)).
					Get("/superuser-check", func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNoContent)
					})
			})
		})
	})

	return r
}
