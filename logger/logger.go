package logger

import (
	"github.com/hibiken/asynq"
	"net/http"
)

type ILogger interface {
	Info(message string)
	Error(message string)
	RequestLogger() func(next http.Handler) http.Handler
	TaskLogger() func(handler asynq.Handler) asynq.Handler
}
