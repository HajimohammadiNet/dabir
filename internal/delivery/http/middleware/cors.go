package middleware

import (
	"net/http"
	"strings"

	"github.com/hajimohammadinet/dabir/internal/config"
)

func CORS(cfg config.AppConfig) func(http.Handler) http.Handler {
	allowedOrigins := splitAndTrim(cfg.CORSAllowedOrigins)
	allowedMethods := cfg.CORSAllowedMethods
	allowedHeaders := cfg.CORSAllowedHeaders

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if isOriginAllowed(origin, allowedOrigins) {
				if contains(allowedOrigins, "*") {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Max-Age", "300")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if contains(allowedOrigins, "*") {
		return true
	}

	if origin == "" {
		return false
	}

	return contains(allowedOrigins, origin)
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
