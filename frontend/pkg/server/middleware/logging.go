package middleware

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)

		lrw := &loggingResponseWriter{rw, http.StatusOK}
		next.ServeHTTP(lrw, r)

		log.Printf("[%s] %s - %d", r.Method, r.URL.Path, lrw.statusCode)
	})
}
