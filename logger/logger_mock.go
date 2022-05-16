package logger

import (
	"context"
	"github.com/hibiken/asynq"
	"net/http"
)

type mockLogger struct{}

func NewMockLogger() mockLogger {
	return mockLogger{}
}

func (mockLogger) Info(message string) {}

func (mockLogger) Error(message string) {}

func (mockLogger) RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (mockLogger) TaskLogger() func(handler asynq.Handler) asynq.Handler {
	return func(next asynq.Handler) asynq.Handler {
		fn := func(ctx context.Context, task *asynq.Task) error {
			err := next.ProcessTask(ctx, task)
			return err
		}
		return asynq.HandlerFunc(fn)
	}
}
