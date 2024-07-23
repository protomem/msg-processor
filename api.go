package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type APIServerOptions struct {
	ListenAddr string
	BaseURL    string
}

type APIServer struct {
	opts APIServerOptions
	log  *slog.Logger
	srv  *http.Server
}

func NewAPIServer(log *slog.Logger, opts APIServerOptions) *APIServer {
	return &APIServer{
		opts: opts,
		log:  log.With("component", "apiServer"),
		srv:  &http.Server{Addr: opts.ListenAddr},
	}
}

func (s *APIServer) ListenAddr() string {
	return s.srv.Addr
}

func (s *APIServer) Run() error {
	h := s.setupRoutes()
	s.srv.Handler = h

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *APIServer) setupRoutes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/health", makeHTTPHandleFunc(s.handleHealth))

	return router
}

func (s *APIServer) handleHealth(w http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(w, http.StatusOK, JSONObject{"status": "OK"})
}

type APIHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f APIHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			_ = WriteJSON(w, http.StatusInternalServerError, APIError{Error: err.Error()})
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
