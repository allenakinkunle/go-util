package logger

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *zapLogger {
	return &zapLogger{
		logger: logger,
	}
}

func (z *zapLogger) Info(message string) {
	z.logger.Info(message)
}

func (z *zapLogger) Error(message string) {
	z.logger.Error(message)
}

func (z *zapLogger) RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				z.logger.Info("Served",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.String("proto", r.Proto),
					zap.Int("time", int(time.Since(t1).Milliseconds())),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (z *zapLogger) TaskLogger() func(handler asynq.Handler) asynq.Handler {
	return func(next asynq.Handler) asynq.Handler {
		fn := func(ctx context.Context, task *asynq.Task) error {
			t1 := time.Now()
			defer func() {
				z.logger.Info("Task",
					zap.String("name", task.Type()),
					zap.Int("time", int(time.Since(t1).Milliseconds())),
				)
			}()
			err := next.ProcessTask(ctx, task)
			if err != nil {
				return err
			}
			return nil
		}
		return asynq.HandlerFunc(fn)
	}
}
