package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(srw, r)

		log.Printf("request method=%s path=%s status=%d duration=%s", r.Method, r.URL.Path, srw.status, time.Since(start))
	})
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func LogHandlerError(r *http.Request, message string, err error) {
	if err == nil {
		log.Printf("handler error method=%s path=%s msg=%s", r.Method, r.URL.Path, message)
		return
	}

	log.Printf("handler error method=%s path=%s msg=%s err=%v", r.Method, r.URL.Path, message, err)
}
