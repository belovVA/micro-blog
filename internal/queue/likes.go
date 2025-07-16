package queue

import (
	"context"
	"log/slog"
	"sync"

	"micro-blog/internal/logger"
	"micro-blog/internal/model"
)

type LikeHandler interface {
	HandleLike(ctx context.Context, like *model.Like) error
}

type LikeQueue struct {
	queue   chan *model.Like
	done    chan struct{}
	wg      sync.WaitGroup
	logger  logger.Logger
	handler LikeHandler
}

func NewLikeQueue(serv LikeHandler, sizeBuffer int, log logger.Logger) *LikeQueue {
	if sizeBuffer < 1 {
		sizeBuffer = 1
	}

	q := &LikeQueue{
		queue:   make(chan *model.Like, sizeBuffer),
		handler: serv,
		done:    make(chan struct{}),
		logger:  log,
	}

	q.wg.Add(1)
	go q.worker()

	return q
}

func (q *LikeQueue) Enqueue(like *model.Like) {
	q.queue <- like
}

func (q *LikeQueue) worker() {
	for {
		select {

		case event := <-q.queue:
			err := q.handler.HandleLike(context.Background(), event)
			if err != nil {
				q.logger.Error("failed to like post", slog.String("error", err.Error()))
			}

		case <-q.done:
			close(q.queue)
			for event := range q.queue {
				err := q.handler.HandleLike(context.Background(), event)
				if err != nil {
					q.logger.Error("failed to like post", slog.String("error", err.Error()))
				}
			}
			return
		}
	}

}

func (q *LikeQueue) Close() {
	close(q.done)
	q.wg.Wait()
}
