package main

import (
	"log/slog"
	"net/http"
	"time"
)

func UseMiddleware(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}

	return next
}

func (s *APIServer) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWrapper(w)

		begin := time.Now()
		next.ServeHTTP(rw, r)
		end := time.Now()

		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttr := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("repsonse", "status", rw.StatusCode, "size", rw.BytesCount, "duration", end.Sub(begin))

		s.log.Info("access log", userAttr, requestAttrs, responseAttrs)
	})
}

func (s *APIServer) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.Error("panic occurred", "error", err)
				_ = WriteJSON(w, http.StatusInternalServerError, APIError{Error: http.StatusText(http.StatusInternalServerError)})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWrapper struct {
	StatusCode    int
	BytesCount    int
	headerWritten bool
	wrapped       http.ResponseWriter
}

func newResponseWrapper(wrapped http.ResponseWriter) *responseWrapper {
	return &responseWrapper{
		wrapped: wrapped,
	}
}

func (rw *responseWrapper) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.wrapped.WriteHeader(statusCode)

	if !rw.headerWritten {
		rw.StatusCode = statusCode
		rw.headerWritten = true
	}
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	rw.headerWritten = true

	n, err := rw.wrapped.Write(b)
	rw.BytesCount += n
	return n, err
}

func (rw *responseWrapper) Unwrap() http.ResponseWriter {
	return rw.wrapped
}
