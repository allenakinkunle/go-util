package logger

import (
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
