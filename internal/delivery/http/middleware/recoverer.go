package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/hajimohammadinet/dabir/internal/delivery/http/response"
)

func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error(
						"panic recovered",
						"request_id", GetRequestID(r.Context()),
						"method", r.Method,
						"path", r.URL.Path,
						"panic", recovered,
						"stack", string(debug.Stack()),
					)

					response.Error(
						w,
						http.StatusInternalServerError,
						"INTERNAL_SERVER_ERROR",
						"internal server error",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
