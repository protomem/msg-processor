package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/protomem/msg-processor/pkg/ctxstore"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

func NewScheduler() quartz.Scheduler {
	return quartz.NewStdScheduler()
}

func RunTaskReadProcessingMessages(
	scheduler quartz.Scheduler, baseLog *slog.Logger,
	store Storage, queue Queue,
	runTimeout time.Duration,
) error {
	const taskName = "readProcessingMessages"
	baseLog = baseLog.With("task", taskName)

	task := job.NewFunctionJob(func(ctx context.Context) (struct{}, error) {
		_, log := setupMetadataTask(ctx, baseLog)

		log.Info("starting")
		defer log.Info("finished")

		log.Debug("processing messages")

		return struct{}{}, nil
	})

	return scheduler.ScheduleJob(
		quartz.NewJobDetail(task, quartz.NewJobKey(taskName)),
		quartz.NewSimpleTrigger(runTimeout),
	)
}

func setupMetadataTask(baseCtx context.Context, baseLog *slog.Logger) (ctx context.Context, log *slog.Logger) {
	tid := genTraceId()
	ctx = ctxstore.With(baseCtx, TraceIDKey, tid)
	log = baseLog.With(TraceIDKey.String(), tid)
	return
}
