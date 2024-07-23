package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/protomem/msg-processor/pkg/ctxstore"
)

type APIServerOptions struct {
	ListenAddr string
	BaseURL    string
}

type APIServer struct {
	opts APIServerOptions
	log  *slog.Logger
	srv  *http.Server

	store Storage
	queue Queue
}

func NewAPIServer(log *slog.Logger, store Storage, queue Queue, opts APIServerOptions) *APIServer {
	return &APIServer{
		opts: opts,
		log:  log.With("component", "apiServer"),
		srv:  &http.Server{Addr: opts.ListenAddr},

		store: store,
		queue: queue,
	}
}

func (s *APIServer) ListenAddr() string {
	return s.srv.Addr
}

func (s *APIServer) Run() error {
	h := s.setupRoutes()
	h = UseMiddleware(h, s.traceId, s.logAccess, s.recovery)
	s.srv.Handler = h

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

type APIHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHTTPHandleFunc(log *slog.Logger, handler string, fn APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxstore.With(r.Context(), HandlerKey, handler)
		if err := fn(w, r.WithContext(ctx)); err != nil {
			log.Warn(
				"failed to process request",
				"error", err,
				TraceIDKey.String(), ctxstore.MustFrom[string](ctx, TraceIDKey),
				HandlerKey.String(), ctxstore.MustFrom[string](ctx, HandlerKey),
			)
			_ = WriteJSON(w, http.StatusInternalServerError, APIError{Error: http.StatusText(http.StatusInternalServerError)})
		}
	}
}

type APIError struct {
	Error string `json:"error"`
}

type JSONObject map[string]any

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
