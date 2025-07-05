package handler

import (
	"context"
	"log/slog"
	"net/http"

	"micro-blog/internal/middleware"
)

type Logger interface {
	Info(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
}

type Service interface {
}

type Router struct {
	service Service
}

func NewRouter(service Service, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	//router := &Router{service: service}

	// Инициализация middleware
	validator := middleware.NewValidator().Middleware
	loggerMw := middleware.ContextLoggerMiddleware(logger)
	recovery := middleware.Recovery(logger)

	// Утилита для оборачивания хендлера
	wrap := func(h http.Handler) http.Handler {
		return recovery(loggerMw(validator(h)))
	}

	// Заглушка на "/"
	r.Handle("/", wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Stub endpoint is alive"))
	})))

	return r
}

// Позволяет ограничить вызов ручки только конкретным методом
func methodOnly(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	})
}
