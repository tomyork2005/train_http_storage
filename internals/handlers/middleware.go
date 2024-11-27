package handlers

import (
	"log/slog"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("method", r.Method),
			slog.String("url", r.URL.String())).Info("handle logger")
		next.ServeHTTP(w, r)
	})
}
