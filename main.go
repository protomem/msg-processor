package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/protomem/msg-processor/pkg/env"
)

var _cfgFile = flag.String("cfg", "", "path to config file")

func init() {
	flag.Parse()
}

func main() {
	ctx := context.Background()
	log := NewLogger()

	if *_cfgFile != "" {
		if err := env.Load(*_cfgFile); err != nil {
			log.Error("failed to load config", "error", err)
			panic(err)
		}
	}

	var store Storage
	{
		var opts PgStorageOptions
		opts.DSN = env.GetString("STORE_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
		opts.Ping = env.GetBool("STORE_PING", false)
		opts.Automigrate = env.GetBool("STORE_MIGRATE", false)

		var err error
		store, err = NewPgStorage(ctx, log, opts)
		if err != nil {
			log.Error("failed to create storage", "error", err)
			panic(err)
		}
	}
	defer func() {
		if err := store.Close(ctx); err != nil {
			log.Error("failed to close storage", "error", err)
		}
	}()

	var queue Queue
	{
		var opts KafkaQueueOptions
		opts.Addrs = env.GetString("QUEUE_ADDRS", "localhost:9092")
		opts.Topic = env.GetString("QUEUE_TOPIC", "messages")

		var err error
		queue, err = NewKafkaQueue(ctx, log, opts)
		if err != nil {
			log.Error("failed to create queue")
			panic(err)
		}
	}
	defer func() {
		if err := queue.Close(ctx); err != nil {
			log.Error("failed to close queue")
		}
	}()

	var srv *APIServer
	{
		var opts APIServerOptions
		opts.ListenAddr = env.GetString("LISTEN_ADDR", ":8080")
		opts.BaseURL = env.GetString("BASE_URL", "http://localhost:8080")

		srv = NewAPIServer(log, store, queue, opts)
	}

	scheduler := NewScheduler()
	scheduler.Start(ctx)
	{
		var errs error

		if err := RunTaskReadProcessingMessages(
			scheduler, log,
			store, queue,
			env.GetDuration("READ_PROC_MSGS_INTERVAL", 1*time.Second), env.GetDuration("READ_PROC_MSGS_TIMEOUT", 30*time.Second),
		); err != nil {
			errs = errors.Join(errs, err)
		}

		if errs != nil {
			log.Error("failed to run tasks", "error", errs)
			panic(errs)
		}
	}

	shutdownErrCh := make(chan error)
	go func() {
		<-quit()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)

		scheduler.Stop()
		scheduler.Wait(ctx)

		shutdownErrCh <- err
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
