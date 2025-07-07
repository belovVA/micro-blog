package main

import (
	"context"
	log "log/slog"

	"micro-blog/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Error("failed to initialize app", "error", err)
		panic(err)
	}
	log.Info("starting server")
	err = a.Run()
	if err != nil {
		log.Error("failed to run app", "error", err)
		panic(err)
	}
}
