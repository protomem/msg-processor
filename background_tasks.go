package main

import (
	"context"
	"encoding/json"
	"errors"
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
	runInterval time.Duration, runTimeout time.Duration,
) error {
	const taskName = "readProcessingMessages"
	baseLog = baseLog.With("task", taskName)

	// TODO: Commit messages after they are processed
	task := job.NewFunctionJob(func(ctx context.Context) (struct{}, error) {
		ctx, log := setupMetadataTask(ctx, baseLog)

		ctx, cancel := context.WithTimeout(ctx, runTimeout)
		defer cancel()

		log.Info("starting")
		defer log.Info("finished")

		if count, err := store.CountProcessingMessages(ctx); err != nil {
			log.Error("failed to count messages", "error", err)
			return struct{}{}, err
		} else if count == 0 {
			log.Debug("no messages to process")
			return struct{}{}, nil
		}

		msgsCh := make(chan Message)
		go func() {
			defer close(msgsCh)

			ctx, cancel := context.WithTimeout(ctx, runTimeout/2)
			defer cancel()

			for {
				evt, err := queue.ReadEvent(ctx)
				if err != nil {
					if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
						log.Error("failed to read event", "error", err)
					}
					break
				}

				var msg Message
				if err := json.Unmarshal(evt.Value, &msg); err != nil {
					log.Error("failed to unmarshal message", "error", err)
					break
				}

				msgsCh <- msg
			}
		}()

		msgs := make([]Message, 0)
		for msg := range msgsCh {
			msgs = append(msgs, msg)
		}

		msgIds := make([]uint64, 0, len(msgs))
		for _, msg := range msgs {
			msgIds = append(msgIds, msg.ID)
		}

		log.Debug("processed messages", "countMsgs", len(msgs))

		if err := store.UpdateStatusMessages(ctx, msgIds, MessageCompleted); err != nil {
			log.Error("failed to update messages status", "error", err)
			return struct{}{}, err
		}

		return struct{}{}, nil
	})

	return scheduler.ScheduleJob(
		quartz.NewJobDetail(task, quartz.NewJobKey(taskName)),
		quartz.NewSimpleTrigger(runInterval),
	)
}

func setupMetadataTask(baseCtx context.Context, baseLog *slog.Logger) (ctx context.Context, log *slog.Logger) {
	tid := genTraceId()
	ctx = ctxstore.With(baseCtx, TraceIDKey, tid)
	log = baseLog.With(TraceIDKey.String(), tid)
	return
}
