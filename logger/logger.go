package logger

import "net/http"

type ILogger interface {
	Info(message string)
	Error(message string)
	RequestLogger() func(next http.Handler) http.Handler
}
