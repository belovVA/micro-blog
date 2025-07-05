package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

func ContextLoggerMiddleware(baseLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "logger", baseLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
