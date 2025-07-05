package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"micro-blog/internal/handler/pkg/response"
)

type RecoveryLogger interface {
	Error(msg string, args ...any)
}

func Recovery(logger RecoveryLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("recovered from panic: %v\nstack trace:\n%s", err, debug.Stack()))
					response.WriteError(w, "Something went wrong", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
