package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"micro-blog/internal/handler/pkg/response"
)

func Recovery(baseLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqLogger := baseLogger.With(
				slog.String("request_id", uuid.NewString()),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			defer func() {
				if rec := recover(); rec != nil {
					reqLogger.ErrorContext(
						r.Context(),
						"recovered from panic",
						slog.Any("panic", rec),
						slog.String("stack", string(debug.Stack())),
					)

					response.WriteError(w, "Something went wrong", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
