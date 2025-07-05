package app

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micro-blog/internal/config"
	"micro-blog/internal/config/env"
	"micro-blog/internal/handler"
	"micro-blog/internal/repository"
	"micro-blog/internal/service"
	"micro-blog/pkg/logger"
)

type App struct {
	httpCfg config.HTTPConfig
	router  http.Handler
}

func NewApp(ctx context.Context) (*App, error) {
	logger := logger.InitLogger()

	htppCfg, err := env.HTTPConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading http config: %w", err)
	}
	//init repo
	repo := repository.NewRepository()

	// init service
	serv := service.NewService(repo)

	//init router
	r := handler.NewRouter(serv, logger)

	return &App{
			router:  r,
			httpCfg: htppCfg,
		},
		nil

}

func (a *App) Run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.httpCfg.GetPort()),
		Handler:      a.router,
		ReadTimeout:  a.httpCfg.GetTimeout(),
		WriteTimeout: a.httpCfg.GetTimeout(),
		IdleTimeout:  a.httpCfg.GetIdleTimeout(),
	}

	// Запуск сервера
	go func() {
		log.Info("Starting HTTP server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server ListenAndServe failed", log.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", log.Any("err", err))
		return err
	}

	select {
	case <-ctx.Done():
		log.Warn("Shutdown timeout exceeded")
	default:
		log.Info("Server exited gracefully")
	}

	return nil
}
