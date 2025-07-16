package logger

import (
	"context"
	"log/slog"
	"sync"
)

type Logger interface {
	Info(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	InfoContext(ctx context.Context, msg string, attrs ...slog.Attr)
	ErrorContext(ctx context.Context, msg string, attrs ...slog.Attr)
	With(args ...any) Logger
}

type Event struct {
	Ctx     context.Context
	Level   slog.Level
	Message string
	Attrs   []slog.Attr
}

type AsyncLogger struct {
	logChan   chan Event
	done      chan struct{}
	baseAttrs []slog.Attr
	wg        sync.WaitGroup
}

func NewAsyncLogger(bufferSize int) *AsyncLogger {
	if bufferSize < 1 {
		bufferSize = 1
	}

	al := &AsyncLogger{
		logChan:   make(chan Event, bufferSize),
		done:      make(chan struct{}),
		baseAttrs: nil,
	}
	al.wg.Add(1)
	go al.listen()
	return al
}

func (a *AsyncLogger) listen() {
	defer a.wg.Done()
	for {
		select {
		case event := <-a.logChan:
			allAttrs := append(a.baseAttrs, event.Attrs...)
			slog.LogAttrs(event.Ctx, event.Level, event.Message, allAttrs...)
		case <-a.done:
			close(a.logChan)
			// теперь прочитаем все, что осталось в logChan
			for event := range a.logChan {
				allAttrs := append(a.baseAttrs, event.Attrs...)
				slog.LogAttrs(event.Ctx, event.Level, event.Message, allAttrs...)
			}
			return
		}
	}
}

func (a *AsyncLogger) Info(msg string, attrs ...slog.Attr) {
	a.InfoContext(context.Background(), msg, attrs...)
}
func (a *AsyncLogger) InfoContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	a.logChan <- Event{
		Ctx:     ctx,
		Level:   slog.LevelInfo,
		Message: msg,
		Attrs:   attrs,
	}
}

func (a *AsyncLogger) Error(msg string, attrs ...slog.Attr) {
	a.ErrorContext(context.Background(), msg, attrs...)
}
func (a *AsyncLogger) ErrorContext(ctx context.Context, msg string, attrs ...slog.Attr) {
	a.logChan <- Event{
		Ctx:     ctx,
		Level:   slog.LevelError,
		Message: msg,
		Attrs:   attrs,
	}
}

func (a *AsyncLogger) With(args ...any) Logger {
	// Конвертируем args в []slog.Attr
	newAttrs := convertToAttrs(args...)

	// Создаем новый AsyncLogger с теми же каналами, но с объединенными базовыми атрибутами
	return &AsyncLogger{
		logChan:   a.logChan,
		done:      a.done,
		baseAttrs: append(a.baseAttrs, newAttrs...),
	}
}

func (a *AsyncLogger) Close() {
	close(a.done)
	a.wg.Wait()
}

// convertToAttrs конвертирует произвольный список аргументов в []slog.Attr.
// Можно сделать простой вариант для ключ-значение попарно.
func convertToAttrs(args ...any) []slog.Attr {
	var attrs []slog.Attr
	n := len(args)
	for i := 0; i < n; i += 2 {
		if i+1 >= n {
			// нет пары, можно проигнорировать или создать строковый атрибут с пустым значением
			break
		}
		key, ok := args[i].(string)
		if !ok {
			// если ключ не строка — тоже можно проигнорировать или преобразовать в строку
			key = "unknown"
		}
		val := args[i+1]
		attrs = append(attrs, slog.Any(key, val))
	}
	return attrs
}
