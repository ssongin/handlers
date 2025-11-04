package middleware

import (
	"net/http"
	"time"

	"github.com/ssongin/core"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		var cn string
		if info, ok := GetCertInfo(r); ok && info != nil {
			cn = info.CommonName
		}

		core.GetLogger().Info("Request",
			"method", r.Method,
			"url", r.URL.Path,
			"proto", r.Proto,
			"client_ip", r.RemoteAddr,
			"common_name", cn,
			"status", rw.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
