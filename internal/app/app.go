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
	asyncLogger "micro-blog/internal/logger"
	"micro-blog/internal/queue"
	"micro-blog/internal/repository"
	"micro-blog/internal/service"
	"micro-blog/pkg/pkglogger"
)

type App struct {
	httpCfg   config.HTTPConfig
	router    http.Handler
	logger    *asyncLogger.AsyncLogger
	likeQueue *queue.LikeQueue
}

const (
	bufferLogSize   = 100
	bufferLikeQueue = 100
)

func NewApp(ctx context.Context) (*App, error) {
	_ = pkglogger.InitLogger()
	logger := asyncLogger.NewAsyncLogger(bufferLogSize)

	htppCfg, err := env.HTTPConfigLoad()
	if err != nil {
		return nil, fmt.Errorf("error loading http config: %w", err)
	}

	//init repo
	repo := repository.NewRepository()

	// init service
	serv := service.NewService(repo)

	// init likeQueue
	queueLikes := queue.NewLikeQueue(serv, bufferLikeQueue, logger)

	// ataching queueLike
	serv.PostService.AttachLikeQueue(queueLikes)

	//init router
	r := handler.NewRouter(serv, logger)

	return &App{
			router:    r,
			httpCfg:   htppCfg,
			logger:    logger,
			likeQueue: queueLikes,
		},
		nil

}

func (a *App) Run() error {
	defer a.logger.Close()
	defer a.likeQueue.Close()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.httpCfg.GetPort()),
		Handler:      a.router,
		ReadTimeout:  a.httpCfg.GetTimeout(),
		WriteTimeout: a.httpCfg.GetTimeout(),
		IdleTimeout:  a.httpCfg.GetIdleTimeout(),
	}

	// Запуск сервера
	go func() {
		a.logger.Info("Starting HTTP server", log.Any("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server ListenAndServe failed", log.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	a.logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.ErrorContext(ctx, "Server shutdown failed", log.Any("err", err))
		return err
	}

	select {
	case <-ctx.Done():
		a.logger.Info("Shutdown timeout exceeded")
	default:
		a.logger.Info("Server exited gracefully")
	}

	return nil
}
