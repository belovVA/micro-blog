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

	closeOnce sync.Once
	closed    bool
	mu        sync.Mutex
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
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		q.logger.Info("likeQueue is closed; skipping enqueue")
		return
	}

	q.queue <- like
}

func (q *LikeQueue) worker() {
	defer q.wg.Done()

	for {
		select {
		case event := <-q.queue:
			q.process(event)

		case <-q.done:
			for {
				select {
				case event := <-q.queue:
					q.process(event)
				default:
					return
				}
			}
		}
	}
}

func (q *LikeQueue) process(like *model.Like) {
	err := q.handler.HandleLike(context.Background(), like)
	if err != nil {
		q.logger.Error("failed to like post", slog.String("error", err.Error()))
	}
}

func (q *LikeQueue) Close() {
	q.closeOnce.Do(func() {
		q.mu.Lock()
		q.closed = true
		q.mu.Unlock()

		close(q.done)
		q.wg.Wait()
	})
}
