package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := NewLogger()

	var srv *APIServer
	{
		var opts APIServerOptions
		opts.ListenAddr = ":8080"
		opts.BaseURL = "http://localhost:8080"

		srv = NewAPIServer(log, opts)
	}

	shutdownErrCh := make(chan error)
	go func() {
		<-quit()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		shutdownErrCh <- srv.Shutdown(ctx)
	}()

	log.Info("starting server", "addr", srv.ListenAddr())
	defer log.Info("shutting down server")

	if err := srv.Run(); err != nil {
		log.Error("failed to run server", "error", err)
		panic(err)
	}

	if err := <-shutdownErrCh; err != nil {
		log.Error("failed to shutdown server", "error", err)
		panic(err)
	}
}

func quit() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	return ch
}
