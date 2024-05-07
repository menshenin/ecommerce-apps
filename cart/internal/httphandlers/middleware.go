package httphandlers

import (
	"log/slog"
	"net/http"
)

// HTTPMiddleware Middleware
type HTTPMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

// WithLog Логирующая Middleware
func WithLog(logger *slog.Logger) HTTPMiddleware {
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			logger.InfoContext(request.Context(), "request received", "path", request.URL.String())
			handlerFunc(writer, request)
		}
	}
}
