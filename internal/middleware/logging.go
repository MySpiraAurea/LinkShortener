package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

        next.ServeHTTP(ww, r)

        duration := time.Since(start)

        slog.Info("HTTP запрос",
            "method", r.Method,
            "uri", r.URL.RequestURI(),
            "duration", duration.Milliseconds(),
            "status", ww.statusCode,
            "size", ww.size,
        )
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    size       int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    if rw.statusCode == 0 {
        rw.statusCode = http.StatusOK
    }
    size, err := rw.ResponseWriter.Write(b)
    rw.size += size
    return size, err
}